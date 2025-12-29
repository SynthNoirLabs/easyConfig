package config

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"easyConfig/pkg/util/paths"
)

// ProfileItem represents a snapshot of a single config file.
type ProfileItem struct {
	Path     string    `json:"path"`
	Provider string    `json:"provider"`
	Scope    Scope     `json:"scope"`
	Content  string    `json:"content"`
	TakenAt  time.Time `json:"takenAt"`
}

// Profile is a saved set of configuration files.
type Profile struct {
	Name      string        `json:"name"`
	Items     []ProfileItem `json:"items"`
	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt"`
}

// ProfileSummary is a lightweight view for listings.
type ProfileSummary struct {
	Name      string    `json:"name"`
	ItemCount int       `json:"itemCount"`
	UpdatedAt time.Time `json:"updatedAt"`
}

var profileNameRe = regexp.MustCompile(`^[A-Za-z0-9._-]+$`)
var backupRe = regexp.MustCompile(`\.\d{4}-\d{2}-\d{2}T\d{2}-\d{2}-\d{2}Z\.bak$`)

// SaveProfile snapshots all discovered configs into a named profile.
func (s *DiscoveryService) SaveProfile(name, projectPath string) error {
	name = sanitizeProfileName(name)
	if name == "" {
		return fmt.Errorf("invalid profile name")
	}

	root, err := profilesRoot()
	if err != nil {
		return err
	}

	items, err := s.DiscoverAll(context.Background(), projectPath)
	if err != nil {
		return fmt.Errorf("discovering configs: %w", err)
	}

	var snaps []ProfileItem
	for _, item := range items {
		// Skip system-scoped files (users shouldn't modify system configs like /etc/gitconfig)
		// This prevents permission errors when attempting to apply profiles
		if item.Scope == ScopeSystem {
			continue
		}
		content, err := s.ReadConfig(item.Path)
		if err != nil {
			// skip missing/unreadable files
			continue
		}
		snaps = append(snaps, ProfileItem{
			Path:     item.Path,
			Provider: item.Provider,
			Scope:    item.Scope,
			Content:  content,
			TakenAt:  time.Now(),
		})
	}

	prof := Profile{
		Name:      name,
		Items:     snaps,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	data, err := json.MarshalIndent(prof, "", "  ")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(root, 0o700); err != nil {
		return fmt.Errorf("create profiles dir: %w", err)
	}

	path := filepath.Join(root, name+".json")
	return os.WriteFile(path, data, 0o600)
}

// ListProfiles lists saved profiles.
func (s *DiscoveryService) ListProfiles() ([]ProfileSummary, error) {
	root, err := profilesRoot()
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(root)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	summaries := []ProfileSummary{}
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != jsonExt {
			continue
		}
		//nolint:gosec // G304: Path constructed from user input (profile name)
		data, err := os.ReadFile(filepath.Join(root, e.Name()))
		if err != nil {
			continue
		}
		var prof Profile
		if err := json.Unmarshal(data, &prof); err != nil {
			continue
		}
		summaries = append(summaries, ProfileSummary{
			Name:      prof.Name,
			ItemCount: len(prof.Items),
			UpdatedAt: prof.UpdatedAt,
		})
	}

	return summaries, nil
}

// ApplyProfile writes the saved snapshot back to disk.
// Returns the list of written file paths for reporting.
func (s *DiscoveryService) ApplyProfile(name string) ([]string, error) {
	prof, err := s.loadProfile(name)
	if err != nil {
		return nil, err
	}

	var written []string

	for _, item := range prof.Items {
		// Create backup before overwriting
		if err := createBackup(item.Path); err != nil {
			return written, fmt.Errorf("backup %s: %w", item.Path, err)
		}

		if err := os.MkdirAll(filepath.Dir(item.Path), 0o750); err != nil {
			return written, fmt.Errorf("create dirs for %s: %w", item.Path, err)
		}
		if err := s.SaveConfig(item.Path, item.Content); err != nil {
			return written, fmt.Errorf("write %s: %w", item.Path, err)
		}
		written = append(written, item.Path)
	}
	return written, nil
}

// Backup represents a single backup file.
type Backup struct {
	Path      string    `json:"path"`
	Timestamp time.Time `json:"timestamp"`
}

// ListBackups lists available backups for a given original file path.
func (s *DiscoveryService) ListBackups(originalPath string) ([]Backup, error) {
	dir := filepath.Dir(originalPath)
	base := filepath.Base(originalPath)
	pattern := filepath.Join(dir, base+".*.bak")

	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	var backups []Backup
	for _, match := range matches {
		// Extract timestamp from filename like `config.json.2023-10-27T10-30-00Z.bak`
		parts := strings.Split(match, ".")
		if len(parts) < 3 {
			continue
		}
		// Assuming format is base.timestamp.bak
		tsPart := parts[len(parts)-2]
		ts, err := time.Parse("2006-01-02T15-04-05Z", tsPart)
		if err != nil {
			continue
		}
		backups = append(backups, Backup{Path: match, Timestamp: ts})
	}

	// Sort by most recent first
	sort.Slice(backups, func(i, j int) bool {
		return backups[i].Timestamp.After(backups[j].Timestamp)
	})

	return backups, nil
}

// RestoreBackup restores a specific backup file.
func (s *DiscoveryService) RestoreBackup(backupPath string) error {
	originalPath := backupRe.ReplaceAllString(backupPath, "")
	if originalPath == backupPath { // No match
		return fmt.Errorf("could not determine original path from backup: invalid format")
	}

	content, err := os.ReadFile(backupPath)
	if err != nil {
		return err
	}
	return os.WriteFile(originalPath, content, 0o600)
}

func createBackup(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No backup needed for new files
		}
		return err
	}

	// Using a timestamp for unique backup names
	timestamp := time.Now().UTC().Format("2006-01-02T15-04-05Z")
	backupPath := fmt.Sprintf("%s.%s.bak", filePath, timestamp)

	if err := os.WriteFile(backupPath, content, 0o600); err != nil {
		return err
	}
	// This assumes that the DiscoveryService logger is accessible here.
	// If not, we might need to pass it down from ApplyProfile.
	// For now, let's assume we can't easily access it and will add it in a future step if needed.
	return cleanupBackups(filePath, 3, nil)
}

func cleanupBackups(originalPath string, keep int, logger *slog.Logger) error {
	dir := filepath.Dir(originalPath)
	base := filepath.Base(originalPath)
	pattern := filepath.Join(dir, base+".*.bak")

	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	if len(matches) <= keep {
		return nil
	}

	// Sort by name (timestamp)
	sort.Strings(matches)

	// Remove the oldest ones
	for i := 0; i < len(matches)-keep; i++ {
		if err := os.Remove(matches[i]); err != nil {
			if logger != nil {
				logger.Warn("failed to remove old backup", "path", matches[i], "err", err)
			}
		}
	}
	return nil
}

// ConfigChange represents a pending change to a file from a profile.
type ConfigChange struct {
	Path       string `json:"path"`
	Status     string `json:"status"` // "modified", "added", "removed" (for future)
	NewContent string `json:"newContent"`
	Content    string `json:"content"`
}

// PreviewApplyProfile shows the changes that would be made by applying a profile.
func (s *DiscoveryService) PreviewApplyProfile(name string) ([]ConfigChange, error) {
	prof, err := s.loadProfile(name)
	if err != nil {
		return nil, err
	}

	var changes []ConfigChange
	for _, item := range prof.Items {
		currentContent, err := os.ReadFile(item.Path)
		status := "modified"
		if err != nil {
			if os.IsNotExist(err) {
				status = "added"
				currentContent = []byte{}
			} else {
				continue // Skip files we can't read
			}
		}

		newContent, err := calculateDiff(string(currentContent), item.Content)
		if err != nil {
			newContent = "Could not calculate diff"
		}

		changes = append(changes, ConfigChange{
			Path:       item.Path,
			Status:     status,
			NewContent: newContent,
			Content:    item.Content,
		})
	}
	return changes, nil
}

func calculateDiff(before, after string) (string, error) {
	// For simplicity, we'll just show the new content.
	// A proper implementation would use a diff library.
	// Example using go-diff:
	// dmp := diffmatchpatch.New()
	// diffs := dmp.DiffMain(before, after, true)
	// return dmp.DiffPrettyText(diffs), nil
	return after, nil
}

// DeleteProfile removes a saved profile.
func (s *DiscoveryService) DeleteProfile(name string) error {
	name = sanitizeProfileName(name)
	root, err := profilesRoot()
	if err != nil {
		return err
	}
	return os.Remove(filepath.Join(root, name+".json"))
}

func (s *DiscoveryService) loadProfile(name string) (*Profile, error) {
	name = sanitizeProfileName(name)
	root, err := profilesRoot()
	if err != nil {
		return nil, err
	}
	//nolint:gosec // G304: Path constructed from user input (profile name)
	data, err := os.ReadFile(filepath.Join(root, name+".json"))
	if err != nil {
		return nil, fmt.Errorf("profile not found: %w", err)
	}
	var prof Profile
	if err := json.Unmarshal(data, &prof); err != nil {
		return nil, err
	}
	return &prof, nil
}

func profilesRoot() (string, error) {
	base := paths.GetConfigDir("easyconfig")
	if base == "" {
		return "", fmt.Errorf("could not determine config dir")
	}
	return filepath.Join(base, "profiles"), nil
}

func sanitizeProfileName(name string) string {
	name = strings.TrimSpace(name)
	if !profileNameRe.MatchString(name) {
		return ""
	}
	return name
}
