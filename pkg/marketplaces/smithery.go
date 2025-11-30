package marketplaces

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// SmitheryClient handles interactions with the Smithery marketplace
type SmitheryClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewSmitheryClient creates a new SmitheryClient
func NewSmitheryClient() *SmitheryClient {
	return &SmitheryClient{
		BaseURL: "https://api.smithery.ai/v1",
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// FetchPopularServers fetches popular MCP servers from Smithery
func (c *SmitheryClient) FetchPopularServers() ([]MCPPackage, error) {
	client := c.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}

	baseURL := c.BaseURL
	if baseURL == "" {
		baseURL = "https://api.smithery.ai/v1"
	}

	resp, err := client.Get(fmt.Sprintf("%s/packages/popular", baseURL))
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fallbackSmitheryPackages(), nil
	}

	var packages []MCPPackage
	if err := json.NewDecoder(resp.Body).Decode(&packages); err != nil {
		return nil, err
	}

	return packages, nil
}

func fallbackSmitheryPackages() []MCPPackage {
	return []MCPPackage{
		{
			Name:        "@modelcontextprotocol/server-filesystem",
			Description: "Official filesystem server for MCP",
			Vendor:      "Anthropic",
			Source:      "smithery",
			RepoURL:     "https://github.com/modelcontextprotocol/server-filesystem",
			License:     "MIT",
			Verified:    true,
		},
		{
			Name:        "@modelcontextprotocol/server-git",
			Description: "Official Git server for MCP",
			Vendor:      "Anthropic",
			Source:      "smithery",
			RepoURL:     "https://github.com/modelcontextprotocol/server-git",
			License:     "MIT",
			Verified:    true,
		},
		{
			Name:        "@modelcontextprotocol/server-memory",
			Description: "Server for persistent memory",
			Vendor:      "Anthropic",
			Source:      "smithery",
			RepoURL:     "https://github.com/modelcontextprotocol/server-memory",
			License:     "MIT",
			Verified:    true,
		},
		{
			Name:        "mcp-server-postgres",
			Description: "PostgreSQL interface for MCP",
			Vendor:      "Community",
			Source:      "smithery",
			RepoURL:     "https://github.com/mcp-community/mcp-server-postgres",
			License:     "Apache-2.0",
			Verified:    false,
		},
	}
}
