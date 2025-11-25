package marketplaces

// MCPPackage represents an MCP server package
type MCPPackage struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Vendor      string `json:"vendor,omitempty"`
	Source      string `json:"source"` // "smithery", "awesome-punkpeye", "awesome-wong2"
	URL         string `json:"url,omitempty"`
}
