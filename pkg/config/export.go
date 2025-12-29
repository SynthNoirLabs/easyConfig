package config

import (
	"encoding/json"
	"fmt"
	"time"
)

const ExportVersion = "1.0"

// ExportedConfig represents a single configuration file within an exported profile.
type ExportedConfig struct {
	Provider string `json:"provider"`
	Scope    Scope  `json:"scope"`
	Content  string `json:"content"`
}

// ExportedProfile represents a profile prepared for export.
type ExportedProfile struct {
	Name        string           `json:"name"`
	Description string           `json:"description,omitempty"`
	Configs     []ExportedConfig `json:"configs"`
}

// ExportData is the top-level structure for the exported JSON file.
type ExportData struct {
	Version    string            `json:"version"`
	ExportedAt time.Time         `json:"exportedAt"`
	Profiles   []ExportedProfile `json:"profiles"`
}

// ExportProfiles packages specified profiles into a JSON structure for export.
func (s *DiscoveryService) ExportProfiles(names []string) ([]byte, error) {
	var exportedProfiles []ExportedProfile

	for _, name := range names {
		profile, err := s.loadProfile(name)
		if err != nil {
			return nil, fmt.Errorf("failed to load profile %s: %w", name, err)
		}

		var exportedConfigs []ExportedConfig
		for _, item := range profile.Items {
			exportedConfigs = append(exportedConfigs, ExportedConfig{
				Provider: item.Provider,
				Scope:    item.Scope,
				Content:  item.Content,
			})
		}

		exportedProfiles = append(exportedProfiles, ExportedProfile{
			Name:    profile.Name,
			// Description is not part of the current Profile struct.
			Description: "",
			Configs:     exportedConfigs,
		})
	}

	exportData := ExportData{
		Version:    ExportVersion,
		ExportedAt: time.Now(),
		Profiles:   exportedProfiles,
	}

	return json.MarshalIndent(exportData, "", "  ")
}

// ExportAllProfiles packages all available profiles for export.
func (s *DiscoveryService) ExportAllProfiles() ([]byte, error) {
	summaries, err := s.ListProfiles()
	if err != nil {
		return nil, fmt.Errorf("failed to list profiles: %w", err)
	}

	var names []string
	for _, summary := range summaries {
		names = append(names, summary.Name)
	}

	return s.ExportProfiles(names)
}
