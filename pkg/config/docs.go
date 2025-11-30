package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// DocsPage represents a single documentation page stored locally.
// Files live under docs/vendor/<provider>/<date>/<slug>.md or .html.
type DocsPage struct {
	Provider    string `json:"provider"`
	Title       string `json:"title"`
	Slug        string `json:"slug"`
	Date        string `json:"date"`
	HasMarkdown bool   `json:"hasMarkdown"`
	HasHTML     bool   `json:"hasHtml"`
}

// DocsProvider is a collection of pages for a given provider and snapshot date.
type DocsProvider struct {
	Provider string     `json:"provider"`
	Date     string     `json:"date"`
	Pages    []DocsPage `json:"pages"`
}

// ListDocsFromRoot scans docs/vendor under the given root directory and
// returns the providers and their pages using the "latest" snapshot.
func ListDocsFromRoot(root string) ([]DocsProvider, error) {
	base := filepath.Join(root, "docs", "vendor")
	entries, err := os.ReadDir(base)
	if err != nil {
		if os.IsNotExist(err) {
			return []DocsProvider{}, nil
		}
		return nil, fmt.Errorf("read docs dir: %w", err)
	}

	var providers []DocsProvider

	const (
		mdExt   = ".md"
		htmlExt = ".html"
	)

	for _, dir := range entries {
		if !dir.IsDir() {
			continue
		}
		provider := dir.Name()
		latestPath := filepath.Join(base, provider, "latest")
		target, err := filepath.EvalSymlinks(latestPath)
		if err != nil {
			target = latestPath
		}
		info, err := os.Stat(target)
		if err != nil || !info.IsDir() {
			continue
		}
		date := filepath.Base(target)

		files, err := os.ReadDir(target)
		if err != nil {
			continue
		}

		pageMap := make(map[string]*DocsPage)

		for _, f := range files {
			if f.IsDir() {
				continue
			}
			name := f.Name()
			if name == "_sources.txt" || name == ".gitkeep" {
				continue
			}
			ext := filepath.Ext(name)
			if ext != mdExt && ext != htmlExt {
				continue
			}
			slug := name[:len(name)-len(ext)]
			page, ok := pageMap[slug]
			if !ok {
				page = &DocsPage{
					Provider: provider,
					Title:    slug,
					Slug:     slug,
					Date:     date,
				}
				pageMap[slug] = page
			}
			if ext == mdExt {
				page.HasMarkdown = true
			} else if ext == htmlExt {
				page.HasHTML = true
			}
		}

		if len(pageMap) == 0 {
			continue
		}

		dp := DocsProvider{
			Provider: provider,
			Date:     date,
		}
		for _, p := range pageMap {
			dp.Pages = append(dp.Pages, *p)
		}
		providers = append(providers, dp)
	}

	return providers, nil
}

// ReadDocFromRoot returns the contents of a local doc page. It looks under
// docs/vendor/<provider>/latest and tries the requested format first, then
// falls back to the other.
func ReadDocFromRoot(root, provider, slug, format string) (string, error) {
	base := filepath.Join(root, "docs", "vendor", provider, "latest")
	if target, err := filepath.EvalSymlinks(base); err == nil {
		base = target
	}

	var tryExts []string
	switch format {
	case "html":
		tryExts = []string{".html", ".md"}
	default:
		tryExts = []string{".md", ".html"}
	}

	for _, ext := range tryExts {
		path := filepath.Join(base, slug+ext)
		if !FileExists(path) {
			continue
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return "", fmt.Errorf("read doc: %w", err)
		}
		return string(data), nil
	}

	return "", fmt.Errorf("doc not found for provider=%s slug=%s", provider, slug)
}
