package marketplaces

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// AwesomeClient handles fetching servers from Awesome MCP lists
type AwesomeClient struct {
	httpClient *http.Client
}

// NewAwesomeClient creates a new AwesomeClient
func NewAwesomeClient() *AwesomeClient {
	return &AwesomeClient{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// FetchServers fetches servers from supported Awesome lists
func (c *AwesomeClient) FetchServers() ([]MCPPackage, error) {
	var allPackages []MCPPackage

	// Fetch from punkpeye/awesome-mcp-servers
	punkpeyePkgs, err := c.fetchAndParse("https://raw.githubusercontent.com/punkpeye/awesome-mcp-servers/main/README.md", "awesome-punkpeye")
	if err == nil {
		allPackages = append(allPackages, punkpeyePkgs...)
	} else {
		fmt.Printf("Error fetching punkpeye list: %v\n", err)
	}

	// Fetch from wong2/awesome-mcp-servers
	wong2Pkgs, err := c.fetchAndParse("https://raw.githubusercontent.com/wong2/awesome-mcp-servers/main/README.md", "awesome-wong2")
	if err == nil {
		allPackages = append(allPackages, wong2Pkgs...)
	} else {
		fmt.Printf("Error fetching wong2 list: %v\n", err)
	}

	return allPackages, nil
}

func (c *AwesomeClient) fetchAndParse(url, source string) ([]MCPPackage, error) {
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch URL %s: status %d", url, resp.StatusCode)
	}

	return c.parseMarkdown(resp.Body, source)
}

// parseMarkdown parses the README markdown to extract server info
// This is a heuristic parser based on common awesome list formats
func (c *AwesomeClient) parseMarkdown(r io.Reader, source string) ([]MCPPackage, error) {
	var packages []MCPPackage
	scanner := bufio.NewScanner(r)

	// Regex to match markdown links: - [Name](URL) - Description
	// or - [Name](URL) Description
	linkRegex := regexp.MustCompile(`^-\s+\[([^\]]+)\]\(([^)]+)\)\s*-?\s*(.*)$`)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		matches := linkRegex.FindStringSubmatch(line)
		if len(matches) == 4 {
			name := matches[1]
			url := matches[2]
			desc := matches[3]

			// Filter out non-server links (heuristics)
			if strings.Contains(strings.ToLower(name), "awesome") || strings.Contains(url, "#") {
				continue
			}

			// Clean up description
			desc = strings.TrimPrefix(desc, "- ")
			desc = strings.TrimSpace(desc)

			pkg := MCPPackage{
				Name:        name,
				Description: desc,
				URL:         url,
				Source:      source,
				Vendor:      "Community", // Default vendor
			}
			packages = append(packages, pkg)
		}
	}

	return packages, scanner.Err()
}
