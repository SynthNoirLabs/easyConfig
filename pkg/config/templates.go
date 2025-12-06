package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type Template struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Provider    string            `json:"provider"`
	Description string            `json:"description"`
	Files       map[string]FileTemplate `json:"files"`
}

type FileTemplate struct {
	Format  string      `json:"format"`
	Content interface{} `json:"content"`
}

func GetTemplates() ([]Template, error) {
	templates := []Template{}

	// For now, we'll read from the `pkg/config/templates` directory.
	// In a real application, this might be embedded or come from a server.
	templatesDir := "pkg/config/templates"

	files, err := os.ReadDir(templatesDir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			templatePath := filepath.Join(templatesDir, file.Name())
			data, err := os.ReadFile(templatePath)
			if err != nil {
				// Log the error but continue
				continue
			}

			var tpl Template
			if err := json.Unmarshal(data, &tpl); err != nil {
				// Log the error but continue
				continue
			}

			// The ID is the filename without the extension
			tpl.ID = strings.TrimSuffix(file.Name(), ".json")
			templates = append(templates, tpl)
		}
	}

	return templates, nil
}
