package schema

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
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
	if err := os.MkdirAll(outputDir, 0o750); err != nil {
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
	parsed, err := validateHTTPSURL(url)
	if err != nil {
		return err
	}
	if filepath.Base(filename) != filename {
		return fmt.Errorf("invalid filename: must not contain path separators")
	}

	//nolint:gosec // G107: URL is validated and comes from the internal registry.
	resp, err := f.client.Get(parsed.String())
	if err != nil {
		return fmt.Errorf("failed to fetch schema: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch schema, status: %s", resp.Status)
	}

	outputPath := filepath.Join(outputDir, filename)

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(outputPath), 0o750); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create file
	//nolint:gosec // G304: outputPath is created under a caller-chosen outputDir with a validated base filename.
	file, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write content: %w", err)
	}

	return nil
}

func validateHTTPSURL(raw string) (*url.URL, error) {
	parsed, err := url.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("invalid url: %w", err)
	}
	if parsed.Scheme != "https" {
		return nil, fmt.Errorf("invalid url: scheme must be https")
	}
	if parsed.Host == "" {
		return nil, fmt.Errorf("invalid url: missing host")
	}
	return parsed, nil
}
