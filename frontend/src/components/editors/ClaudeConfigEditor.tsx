import type React from "react";
import { useEffect, useState } from "react";
import "./EditorStyles.css";

interface ClaudeConfig {
  globalShortcut?: string;
  theme?: "light" | "dark" | "auto";
  allowBrowser?: boolean;
  allowShell?: boolean;
}

interface ClaudeConfigEditorProps {
  content: string;
  onChange: (newContent: string) => void;
}

const ClaudeConfigEditor: React.FC<ClaudeConfigEditorProps> = ({
  content,
  onChange,
}) => {
  const [config, setConfig] = useState<ClaudeConfig>({});
  const [parseError, setParseError] = useState<string | null>(null);

  useEffect(() => {
    try {
      const parsed = JSON.parse(content || "{}");
      setConfig(parsed);
      setParseError(null);
    } catch (_e) {
      setParseError("Invalid JSON content. Please switch to Code view to fix.");
    }
  }, [content]);

  const updateConfig = (updates: Partial<ClaudeConfig>) => {
    const newConfig = { ...config, ...updates };
    setConfig(newConfig);
    onChange(JSON.stringify(newConfig, null, 2));
  };

  if (parseError) {
    return <div className="editor-error">{parseError}</div>;
  }

  return (
    <div className="form-editor">
      <div className="form-group">
        <label htmlFor="claude-global-shortcut">Global Shortcut</label>
        <input
          id="claude-global-shortcut"
          type="text"
          value={config.globalShortcut || ""}
          onChange={(e) => updateConfig({ globalShortcut: e.target.value })}
          placeholder="e.g. Ctrl+Space"
        />
        <small>Key combination to toggle Claude Desktop</small>
      </div>

      <div className="form-group">
        <label htmlFor="claude-theme">Theme</label>
        <select
          id="claude-theme"
          value={config.theme || "auto"}
          onChange={(e) =>
            updateConfig({ theme: e.target.value as ClaudeConfig["theme"] })
          }
        >
          <option value="auto">Auto (System)</option>
          <option value="light">Light</option>
          <option value="dark">Dark</option>
        </select>
      </div>

      <div className="form-group checkbox-group">
        <label>
          <input
            type="checkbox"
            checked={config.allowBrowser || false}
            onChange={(e) => updateConfig({ allowBrowser: e.target.checked })}
          />
          Allow Browser Tool
        </label>
        <small>Enable Claude to use the built-in browser</small>
      </div>

      <div className="form-group checkbox-group">
        <label>
          <input
            type="checkbox"
            checked={config.allowShell || false}
            onChange={(e) => updateConfig({ allowShell: e.target.checked })}
          />
          Allow Shell Tool
        </label>
        <small>Enable Claude to execute shell commands</small>
      </div>
    </div>
  );
};

export default ClaudeConfigEditor;
