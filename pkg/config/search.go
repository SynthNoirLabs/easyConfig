package config

import (
	"bufio"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"
)

// SearchResult represents the result of a search for a single config item
type SearchResult struct {
	ConfigItem Item    `json:"configItem"`
	Matches    []Match `json:"matches"`
}

// Match represents a single match within a file
type Match struct {
	Line    int    `json:"line"`
	Column  int    `json:"column"`
	Text    string `json:"text"`
	Context string `json:"context"` // Surrounding text
}

// SearchOptions defines the options for a search
type SearchOptions struct {
	CaseSensitive bool `json:"caseSensitive"`
	Regex         bool `json:"regex"`
	WholeWord     bool `json:"wholeWord"`
}

// SearchAll searches across all discovered configuration files for a given query.
func (s *DiscoveryService) SearchAll(query string, options SearchOptions) ([]SearchResult, error) {
	var results []SearchResult
	items, err := s.DiscoverAll("") // Assuming projectPath is not needed for search
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		if !item.Exists {
			continue
		}
		matches, err := findMatches(item.Path, query, options)
		if err != nil {
			// Log the error but continue to the next file
			s.logger.Error("error searching file", "path", item.Path, "error", err)
			continue
		}
		if len(matches) > 0 {
			results = append(results, SearchResult{
				ConfigItem: item,
				Matches:    matches,
			})
		}
	}
	return results, nil
}

func findMatches(filePath, query string, options SearchOptions) ([]Match, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var matches []Match
	scanner := bufio.NewScanner(strings.NewReader(string(content)))
	lineNumber := 0

	var re *regexp.Regexp
	if options.Regex {
		if !options.CaseSensitive {
			query = "(?i)" + query
		}
		re, err = regexp.Compile(query)
		if err != nil {
			return nil, err
		}
	}

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()

		if options.Regex {
			locs := re.FindAllStringIndex(line, -1)
			for _, loc := range locs {
				start, end := loc[0], loc[1]
				col := utf8.RuneCountInString(line[:start])
				matches = append(matches, Match{
					Line:    lineNumber,
					Column:  col + 1,
					Text:    line[start:end],
					Context: line,
				})
			}
		} else {
			searchText := query
			lineText := line
			if !options.CaseSensitive {
				searchText = strings.ToLower(searchText)
				lineText = strings.ToLower(lineText)
			}

			offset := 0
			for {
				idx := strings.Index(lineText[offset:], searchText)
				if idx == -1 {
					break
				}

				start := offset + idx
				end := start + len(searchText)

				if options.WholeWord {
					if start > 0 && isWordChar(lineText[start-1]) {
						offset = end
						continue
					}
					if end < len(lineText) && isWordChar(lineText[end]) {
						offset = end
						continue
					}
				}

				col := utf8.RuneCountInString(line[:start])
				matches = append(matches, Match{
					Line:    lineNumber,
					Column:  col + 1,
					Text:    line[start:end],
					Context: line,
				})
				offset = end
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return matches, nil
}

func isWordChar(r byte) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_'
}
