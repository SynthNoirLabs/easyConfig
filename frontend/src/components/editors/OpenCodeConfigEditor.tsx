import type React from "react";
import { useEffect, useState } from "react";
import "./EditorStyles.css";

interface OpenCodeConfig {
  defaultModel?: string;
  maxTokens?: number;
  temperature?: number;
  apiKey?: string;
}

interface OpenCodeConfigEditorProps {
  content: string;
  onChange: (newContent: string) => void;
}

const OpenCodeConfigEditor: React.FC<OpenCodeConfigEditorProps> = ({
  content,
  onChange,
}) => {
  const [config, setConfig] = useState<OpenCodeConfig>({});
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

  const updateConfig = (updates: Partial<OpenCodeConfig>) => {
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
        <label htmlFor="opencode-default-model">Default Model</label>
        <select
          id="opencode-default-model"
          value={config.defaultModel || "gpt-4"}
          onChange={(e) => updateConfig({ defaultModel: e.target.value })}
        >
          <option value="gpt-4">GPT-4</option>
          <option value="gpt-3.5-turbo">GPT-3.5 Turbo</option>
          <option value="claude-3-opus">Claude 3 Opus</option>
          <option value="claude-3-sonnet">Claude 3 Sonnet</option>
        </select>
        <small>Select the AI model to use by default</small>
      </div>

      <div className="form-group">
        <label htmlFor="opencode-max-tokens">Max Tokens</label>
        <input
          id="opencode-max-tokens"
          type="number"
          value={config.maxTokens || 4096}
          onChange={(e) =>
            updateConfig({ maxTokens: parseInt(e.target.value, 10) || 0 })
          }
        />
        <small>Maximum number of tokens to generate</small>
      </div>

      <div className="form-group">
        <label htmlFor="opencode-temperature">Temperature (0.0 - 1.0)</label>
        <input
          id="opencode-temperature"
          type="number"
          step="0.1"
          min="0"
          max="1"
          value={config.temperature ?? 0.7}
          onChange={(e) =>
            updateConfig({ temperature: parseFloat(e.target.value) })
          }
        />
        <small>Controls randomness: 0 is deterministic, 1 is creative</small>
      </div>

      <div className="form-group">
        <label htmlFor="opencode-api-key">API Key (Optional override)</label>
        <input
          id="opencode-api-key"
          type="password"
          value={config.apiKey || ""}
          onChange={(e) => updateConfig({ apiKey: e.target.value })}
          placeholder="sk-..."
        />
        <small>Leave empty to use environment variable</small>
      </div>
    </div>
  );
};

export default OpenCodeConfigEditor;
