import React, { useState, useEffect, useCallback } from 'react';
import { Search, X } from 'lucide-react';
import { SearchAll } from '../../wailsjs/go/main/App';
import { config } from '../../wailsjs/go/models';
import useDebounce from '../hooks/useDebounce';
import './CommandPalette.css';

interface CommandPaletteProps {
  onSelect: (path: string, line?: number) => void;
  onClose: () => void;
}

const CommandPalette: React.FC<CommandPaletteProps> = ({ onSelect, onClose }) => {
  const [query, setQuery] = useState('');
  const [results, setResults] = useState<config.SearchResult[]>([]);
  const [isOpen, setIsOpen] = useState(true);
  const debouncedQuery = useDebounce(query, 300);

  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'Escape') {
        setIsOpen(false);
        onClose();
      }
    };
    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [onClose]);

  useEffect(() => {
    if (debouncedQuery) {
      const options: config.SearchOptions = {
        caseSensitive: false,
        regex: false,
        wholeWord: false,
      };
      SearchAll(debouncedQuery, options).then(setResults);
    } else {
      setResults([]);
    }
  }, [debouncedQuery]);

  const handleSelect = (result: config.SearchResult, match?: config.Match) => {
    onSelect(result.configItem.path, match?.line);
    setIsOpen(false);
    onClose();
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
            autoFocus
          />
          <button onClick={() => { setIsOpen(false); onClose(); }} className="close-button">
            <X />
          </button>
        </div>
        <div className="results-list">
          {results.map((result) => (
            <div key={result.configItem.path}>
              <div className="file-header" onClick={() => handleSelect(result)}>
                {result.configItem.fileName}
              </div>
              <ul>
                {result.matches.map((match, index) => (
                  <li key={index} onClick={() => handleSelect(result, match)}>
                    <div className="match-line">
                      <span className="line-number">{match.line}:</span>
                      <span className="match-context">{match.context}</span>
                    </div>
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
