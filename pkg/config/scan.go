package config

import (
	"io/fs"
	"path/filepath"
	"strings"
)

// defaultExclusions mirrors ccexp's fast-scan skip list to avoid huge/vendor dirs.
var defaultExclusions = map[string]struct{}{
	".git":         {},
	"node_modules": {},
	"dist":         {},
	"build":        {},
	".next":        {},
	".turbo":       {},
	".cache":       {},
	".idea":        {},
	".vscode":      {},
	"vendor":       {},
}

// fastWalk walks root with a maxDepth relative to root and applies match to files.
// Returns matched absolute paths. Directories in defaultExclusions are skipped early.
func fastWalk(root string, maxDepth int, match func(path string, d fs.DirEntry) bool) ([]string, error) {
	root = filepath.Clean(root)
	var results []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// Continue past unreadable paths
			return nil
		}

		rel, _ := filepath.Rel(root, path)
		if rel == "." {
			rel = ""
		}
		depth := 0
		if rel != "" {
			depth = len(strings.Split(rel, string(filepath.Separator)))
		}
		if depth > maxDepth {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if d.IsDir() {
			if _, skip := defaultExclusions[d.Name()]; skip {
				return filepath.SkipDir
			}
			return nil
		}

		if match(path, d) {
			results = append(results, path)
		}
		return nil
	})
	return results, err
}
