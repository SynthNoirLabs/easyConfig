import { Search, X } from "lucide-react";
import type React from "react";
import { useEffect, useState } from "react";
// @ts-ignore
import { SearchAll } from "../../wailsjs/go/main/App";
import type { config } from "../../wailsjs/go/models";
import useDebounce from "../hooks/useDebounce";
import "./CommandPalette.css";

// Define missing types locally
interface Match {
  line: number;
  context: string;
}

interface SearchResult {
  configItem: config.Item;
  matches: Match[];
}

interface SearchOptions {
  caseSensitive: boolean;
  regex: boolean;
  wholeWord: boolean;
}

interface CommandPaletteProps {
  onSelect: (path: string, line?: number) => void;
  onClose: () => void;
}

const CommandPalette: React.FC<CommandPaletteProps> = ({
  onSelect,
  onClose,
}) => {
  const [query, setQuery] = useState("");
  const [results, setResults] = useState<SearchResult[]>([]);
  const [isOpen, setIsOpen] = useState(true);
  const debouncedQuery = useDebounce(query, 300);

  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === "Escape") {
        setIsOpen(false);
        onClose();
      }
    };
    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, [onClose]);

  useEffect(() => {
    if (debouncedQuery) {
      const options: SearchOptions = {
        caseSensitive: false,
        regex: false,
        wholeWord: false,
      };
      // Cast to any because the generated types are missing
      (SearchAll as any)(debouncedQuery, options).then((res: any) => setResults(res));
    } else {
      setResults([]);
    }
  }, [debouncedQuery]);

  const handleSelect = (result: SearchResult, match?: Match) => {
    onSelect(result.configItem.path, match?.line);
    setIsOpen(false);
    onClose();
  };

  const handleKeyDownSelect = (
    e: React.KeyboardEvent,
    result: SearchResult,
    match?: Match
  ) => {
    if (e.key === "Enter" || e.key === " ") {
      e.preventDefault();
      handleSelect(result, match);
    }
  };

  if (!isOpen) return null;

  return (
    <div className="command-palette-overlay">
      <div className="command-palette">
        <div className="search-input-container">
          <Search className="search-icon" />
          <input
            type="text"
            placeholder="Search all configurations..."
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            // biome-ignore lint/a11y/noAutofocus: Expected behavior for command palette
            autoFocus
          />
          <button
            type="button"
            onClick={() => {
              setIsOpen(false);
              onClose();
            }}
            className="close-button"
            aria-label="Close"
          >
            <X />
          </button>
        </div>
        <div className="results-list">
          {results.map((result) => (
            <div key={result.configItem.path}>
              <button
                type="button"
                className="file-header"
                onClick={() => handleSelect(result)}
              >
                {result.configItem.fileName}
              </button>
              <ul>
                {result.matches.map((match, index) => (
                  <li key={`${result.configItem.path}-${match.line}-${index}`}>
                    <button
                      type="button"
                      className="match-item"
                      onClick={() => handleSelect(result, match)}
                    >
                      <div className="match-line">
                        <span className="line-number">{match.line}:</span>
                        <span className="match-context">{match.context}</span>
                      </div>
                    </button>
                  </li>
                ))}
              </ul>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default CommandPalette;
