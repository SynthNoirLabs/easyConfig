package schema

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Fetcher handles retrieving schemas from remote sources
type Fetcher struct {
	client *http.Client
}

// NewFetcher creates a new Schema Fetcher
func NewFetcher() *Fetcher {
	return &Fetcher{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// FetchAllSchemas iterates through the registry and downloads available JSON schemas
// It saves them to the specified output directory
func (f *Fetcher) FetchAllSchemas(outputDir string) error {
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	for _, info := range Registry {
		// Only fetch JSON schemas for now
		if info.Type != TypeJSON {
			continue
		}

		fmt.Printf("Fetching schema for %s...\n", info.ToolName)
		if err := f.fetchAndSave(info.URL, outputDir, info.ToolName+".schema.json"); err != nil {
			fmt.Printf("Error fetching schema for %s: %v\n", info.ToolName, err)
			// Continue fetching others even if one fails
		}
	}
	return nil
}

func (f *Fetcher) fetchAndSave(url, outputDir, filename string) error {
	resp, err := f.client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	outputPath := filepath.Join(outputDir, filename)
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write content: %w", err)
	}

	return nil
}
