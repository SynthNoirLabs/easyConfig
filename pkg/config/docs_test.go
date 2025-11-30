package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestListDocsFromRoot_EmptyWhenNoDir(t *testing.T) {
	root := t.TempDir()

	providers, err := ListDocsFromRoot(root)
	if err != nil {
		t.Fatalf("ListDocsFromRoot returned error: %v", err)
	}
	if len(providers) != 0 {
		t.Fatalf("expected no providers, got %d", len(providers))
	}
}

func TestListDocsFromRoot_Basic(t *testing.T) {
	const snapDate = "2025-11-30"

	root := t.TempDir()
	base := filepath.Join(root, "docs", "vendor", "claude")
	latestDir := filepath.Join(base, snapDate)

	if err := os.MkdirAll(latestDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// Create latest symlink
	if err := os.Symlink(snapDate, filepath.Join(base, "latest")); err != nil {
		t.Fatalf("symlink: %v", err)
	}

	// Add one md and one html page with same slug
	if err := os.WriteFile(filepath.Join(latestDir, "settings.md"), []byte("# Settings"), 0o644); err != nil {
		t.Fatalf("write md: %v", err)
	}
	if err := os.WriteFile(filepath.Join(latestDir, "settings.html"), []byte("<h1>Settings</h1>"), 0o644); err != nil {
		t.Fatalf("write html: %v", err)
	}

	// and a sources file that should be ignored
	if err := os.WriteFile(filepath.Join(latestDir, "_sources.txt"), []byte("meta"), 0o644); err != nil {
		t.Fatalf("write sources: %v", err)
	}

	providers, err := ListDocsFromRoot(root)
	if err != nil {
		t.Fatalf("ListDocsFromRoot returned error: %v", err)
	}
	if len(providers) != 1 {
		t.Fatalf("expected 1 provider, got %d", len(providers))
	}
	p := providers[0]
	if p.Provider != "claude" {
		t.Fatalf("expected provider 'claude', got %q", p.Provider)
	}
	if p.Date != snapDate {
		t.Fatalf("expected date %q, got %q", snapDate, p.Date)
	}
	if len(p.Pages) != 1 {
		t.Fatalf("expected 1 page, got %d", len(p.Pages))
	}
	page := p.Pages[0]
	if !page.HasMarkdown || !page.HasHTML {
		t.Fatalf("expected both markdown and html flags true, got md=%v html=%v", page.HasMarkdown, page.HasHTML)
	}
}

func TestReadDocFromRoot_PrefersMarkdown(t *testing.T) {
	root := t.TempDir()
	base := filepath.Join(root, "docs", "vendor", "gemini")
	date := "2025-11-30"
	latestDir := filepath.Join(base, date)

	if err := os.MkdirAll(latestDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.Symlink(date, filepath.Join(base, "latest")); err != nil {
		t.Fatalf("symlink: %v", err)
	}

	if err := os.WriteFile(filepath.Join(latestDir, "page.md"), []byte("from-md"), 0o644); err != nil {
		t.Fatalf("write md: %v", err)
	}
	if err := os.WriteFile(filepath.Join(latestDir, "page.html"), []byte("<p>from-html</p>"), 0o644); err != nil {
		t.Fatalf("write html: %v", err)
	}

	content, err := ReadDocFromRoot(root, "gemini", "page", "md")
	if err != nil {
		t.Fatalf("ReadDocFromRoot returned error: %v", err)
	}
	if content != "from-md" {
		t.Fatalf("expected markdown content, got %q", content)
	}
}

func TestReadDocFromRoot_FallbackToHTML(t *testing.T) {
	root := t.TempDir()
	base := filepath.Join(root, "docs", "vendor", "gemini")
	date := "2025-11-30"
	latestDir := filepath.Join(base, date)

	if err := os.MkdirAll(latestDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.Symlink(date, filepath.Join(base, "latest")); err != nil {
		t.Fatalf("symlink: %v", err)
	}
	if err := os.WriteFile(filepath.Join(latestDir, "page.html"), []byte("<p>only-html</p>"), 0o644); err != nil {
		t.Fatalf("write html: %v", err)
	}

	content, err := ReadDocFromRoot(root, "gemini", "page", "md")
	if err != nil {
		t.Fatalf("ReadDocFromRoot returned error: %v", err)
	}
	if content != "<p>only-html</p>" {
		t.Fatalf("expected html fallback content, got %q", content)
	}
}
