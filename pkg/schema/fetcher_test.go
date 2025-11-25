package schema

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestFetcher_FetchAllSchemas(t *testing.T) {
	// 1. Mock Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"mock": "schema"}`))
	}))
	defer ts.Close()

	// 2. Setup temporary registry with mock URL
	originalRegistry := Registry
	defer func() { Registry = originalRegistry }()

	Registry = []Info{
		{
			ToolName: "MockTool",
			Type:     TypeJSON,
			URL:      ts.URL,
		},
		{
			ToolName: "IgnoredTool",
			Type:     TypeDocs, // Should be skipped
			URL:      ts.URL,
		},
	}

	// 3. Setup Output Directory
	tmpDir := t.TempDir()

	// 4. Run Fetcher
	f := NewFetcher()
	err := f.FetchAllSchemas(tmpDir)
	if err != nil {
		t.Fatalf("FetchAllSchemas failed: %v", err)
	}

	// 5. Verify Files
	expectedFile := filepath.Join(tmpDir, "MockTool.schema.json")
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		t.Errorf("Expected schema file not created: %s", expectedFile)
	}

	ignoredFile := filepath.Join(tmpDir, "IgnoredTool.schema.json")
	if _, err := os.Stat(ignoredFile); !os.IsNotExist(err) {
		t.Errorf("Did not expect schema file for docs type: %s", ignoredFile)
	}
}
