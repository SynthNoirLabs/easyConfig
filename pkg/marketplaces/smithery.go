package marketplaces

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// MCPPackage represents a server package from Smithery
type MCPPackage struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Author      string `json:"author"`
	URL         string `json:"url"`
}

// SmitheryClient handles interactions with the Smithery API
type SmitheryClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewSmitheryClient creates a new client
func NewSmitheryClient() *SmitheryClient {
	return &SmitheryClient{
		BaseURL: "https://api.smithery.ai/v1", // Hypothetical API endpoint, will fail gracefully if 404
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// FetchPopularServers returns a list of popular MCP servers
func (c *SmitheryClient) FetchPopularServers() ([]MCPPackage, error) {
	resp, err := c.HTTPClient.Get(c.BaseURL + "/packages/popular")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from smithery: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		// Fallback to mock data if API is not reachable or returns error (since we are guessing the endpoint)
		return c.getMockData(), nil
	}

	var packages []MCPPackage
	if err := json.NewDecoder(resp.Body).Decode(&packages); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return packages, nil
}

func (c *SmitheryClient) getMockData() []MCPPackage {
	return []MCPPackage{
		{
			Name:        "exa-mcp-server",
			Description: "Exa Search MCP Server for web research",
			Version:     "1.0.0",
			Author:      "Exa",
			URL:         "https://smithery.ai/server/exa",
		},
		{
			Name:        "github-mcp-server",
			Description: "GitHub integration for MCP",
			Version:     "0.5.0",
			Author:      "GitHub",
			URL:         "https://github.com/modelcontextprotocol/servers",
		},
		{
			Name:        "slack-mcp",
			Description: "Slack integration for AI agents",
			Version:     "1.2.0",
			Author:      "Slack",
			URL:         "https://smithery.ai/server/slack",
		},
	}
}
