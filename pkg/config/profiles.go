package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
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

	items, err := s.DiscoverAll(projectPath)
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

// DeleteProfile removes a saved profile.
func (s *DiscoveryService) DeleteProfile(name string) error {
	name = sanitizeProfileName(name)
	root, err := profilesRoot()
	if err != nil {
		return err
	}
	return os.Remove(filepath.Join(root, name+".json"))
}

// GetProfileContent returns the content of a single file from a profile.
func (s *DiscoveryService) GetProfileContent(profileName, filePath string) (string, error) {
	prof, err := s.loadProfile(profileName)
	if err != nil {
		return "", err
	}

	for _, item := range prof.Items {
		if item.Path == filePath {
			return item.Content, nil
		}
	}

	return "", fmt.Errorf("file not found in profile")
}

// ListProfileFiles returns the list of files in a profile.
func (s *DiscoveryService) ListProfileFiles(profileName string) ([]ProfileItem, error) {
	prof, err := s.loadProfile(profileName)
	if err != nil {
		return nil, err
	}
	return prof.Items, nil
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
