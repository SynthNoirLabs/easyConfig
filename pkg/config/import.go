package config

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// ImportStrategy defines how to handle conflicts when importing profiles.
type ImportStrategy string

const (
	// ImportStrategySkip skips importing a profile if one with the same name already exists.
	ImportStrategySkip ImportStrategy = "skip"
	// ImportStrategyRename renames the imported profile by adding a suffix if a conflict occurs.
	ImportStrategyRename ImportStrategy = "rename"
	// ImportStrategyOverwrite overwrites the existing profile with the imported one.
	ImportStrategyOverwrite ImportStrategy = "overwrite"
)

// ImportResultStatus defines the outcome of a single profile import.
type ImportResultStatus string

const (
	ImportStatusSuccess   ImportResultStatus = "success"
	ImportStatusSkipped   ImportResultStatus = "skipped"
	ImportStatusRenamed   ImportResultStatus = "renamed"
	ImportStatusOverwrote ImportResultStatus = "overwrote"
	ImportStatusFailed    ImportResultStatus = "failed"
)

// ImportResult provides feedback on the import of a single profile.
type ImportResult struct {
	Name      string             `json:"name"`
	NewName   string             `json:"newName,omitempty"`
	Status    ImportResultStatus `json:"status"`
	Message   string             `json:"message,omitempty"`
	IsConflict bool               `json:"isConflict"`
}

// ImportProfiles processes the provided data, importing profiles based on the chosen strategy.
func (s *DiscoveryService) ImportProfiles(data []byte, strategy ImportStrategy) ([]ImportResult, error) {
	var importData ExportData
	if err := json.Unmarshal(data, &importData); err != nil {
		return nil, fmt.Errorf("failed to parse import data: %w", err)
	}

	if importData.Version != ExportVersion {
		return nil, fmt.Errorf("unsupported import version: %s", importData.Version)
	}

	existingProfiles, err := s.ListProfiles()
	if err != nil {
		return nil, fmt.Errorf("failed to list existing profiles: %w", err)
	}
	existingMap := make(map[string]struct{})
	for _, p := range existingProfiles {
		existingMap[p.Name] = struct{}{}
	}

	var results []ImportResult
	for _, importedProfile := range importData.Profiles {
		result := ImportResult{Name: importedProfile.Name}
		_, exists := existingMap[importedProfile.Name]

		if exists {
			result.IsConflict = true
			switch strategy {
			case ImportStrategySkip:
				result.Status = ImportStatusSkipped
				result.Message = "Profile already exists."
				results = append(results, result)
				continue
			case ImportStrategyRename:
				newName := fmt.Sprintf("%s-imported-%d", importedProfile.Name, time.Now().Unix())
				importedProfile.Name = newName
				result.NewName = newName
				result.Status = ImportStatusRenamed
			case ImportStrategyOverwrite:
				result.Status = ImportStatusOverwrote
			default:
				result.Status = ImportStatusFailed
				result.Message = fmt.Sprintf("Unknown import strategy: %s", strategy)
				results = append(results, result)
				continue
			}
		} else {
			result.Status = ImportStatusSuccess
		}

		profileToSave := Profile{
			Name:      importedProfile.Name,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		for _, cfg := range importedProfile.Configs {
			// This is a simplification. A real implementation would need to
			// resolve paths, which can be machine-specific. For now, we
			// assume paths are placeholders and focus on content.
			profileToSave.Items = append(profileToSave.Items, ProfileItem{
				Provider: cfg.Provider,
				Scope:    cfg.Scope,
				Content:  cfg.Content,
				TakenAt:  time.Now(),
				Path:     "imported", // Placeholder path
			})
		}

		// Use SaveProfile to write the new or updated profile.
		if err := s.saveProfileToDisk(&profileToSave); err != nil {
			result.Status = ImportStatusFailed
			result.Message = fmt.Sprintf("Failed to save profile: %s", err)
		}
		results = append(results, result)
	}

	return results, nil
}

// ImportProfilesFromFile reads data from a local file and imports it.
func (s *DiscoveryService) ImportProfilesFromFile(path string, strategy ImportStrategy) ([]ImportResult, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read import file: %w", err)
	}
	return s.ImportProfiles(data, strategy)
}

// ImportProfilesFromURL fetches data from a URL and imports it.
func (s *DiscoveryService) ImportProfilesFromURL(url string, strategy ImportStrategy) ([]ImportResult, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch URL: status code %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return s.ImportProfiles(data, strategy)
}

// saveProfileToDisk is a helper to encapsulate the JSON marshaling and file writing.
// It is based on the logic in SaveProfile.
func (s *DiscoveryService) saveProfileToDisk(prof *Profile) error {
	root, err := profilesRoot()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(prof, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal profile: %w", err)
	}

	if err := os.MkdirAll(root, 0o700); err != nil {
		return fmt.Errorf("failed to create profiles directory: %w", err)
	}

	path := root + "/" + prof.Name + ".json"
	return os.WriteFile(path, data, 0o600)
}
