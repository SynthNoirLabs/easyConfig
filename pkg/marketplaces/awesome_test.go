package marketplaces

import (
	"strings"
	"testing"
)

func TestParseMarkdown(t *testing.T) {
	markdown := `
# Awesome MCP Servers

## Official
- [filesystem-server](https://github.com/modelcontextprotocol/servers/tree/main/src/filesystem) - Official filesystem server
- [git-server](https://github.com/modelcontextprotocol/servers/tree/main/src/git) - Official Git server

## Community
- [postgres-mcp](https://github.com/example/postgres-mcp) - A PostgreSQL MCP server
- [weather-mcp](https://github.com/example/weather) - Get weather updates
`

	client := NewAwesomeClient()
	packages, err := client.parseMarkdown(strings.NewReader(markdown), "test-source")
	if err != nil {
		t.Fatalf("parseMarkdown failed: %v", err)
	}

	if len(packages) != 4 {
		t.Errorf("Expected 4 packages, got %d", len(packages))
	}

	expected := []struct {
		Name string
		URL  string
	}{
		{"filesystem-server", "https://github.com/modelcontextprotocol/servers/tree/main/src/filesystem"},
		{"git-server", "https://github.com/modelcontextprotocol/servers/tree/main/src/git"},
		{"postgres-mcp", "https://github.com/example/postgres-mcp"},
		{"weather-mcp", "https://github.com/example/weather"},
	}

	for i, pkg := range packages {
		if pkg.Name != expected[i].Name {
			t.Errorf("Package %d name mismatch: got %s, want %s", i, pkg.Name, expected[i].Name)
		}
		if pkg.URL != expected[i].URL {
			t.Errorf("Package %d URL mismatch: got %s, want %s", i, pkg.URL, expected[i].URL)
		}
		if pkg.Source != "test-source" {
			t.Errorf("Package %d source mismatch: got %s, want test-source", i, pkg.Source)
		}
	}
}
