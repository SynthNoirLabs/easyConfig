package marketplaces

// MCPPackage represents an MCP server package
type MCPPackage struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Vendor      string   `json:"vendor,omitempty"`
	Source      string   `json:"source"` // "smithery", "awesome-punkpeye", "awesome-wong2"
	URL         string   `json:"url,omitempty"`
	Version     string   `json:"version,omitempty"`
	Author      string   `json:"author,omitempty"`
	Stars       int      `json:"stars,omitempty"`
	Downloads   int      `json:"downloads,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	RepoURL     string   `json:"repoUrl,omitempty"`
	License     string   `json:"license,omitempty"`
	Verified    bool     `json:"verified,omitempty"`
	Checksum    string   `json:"checksum,omitempty"`
}
