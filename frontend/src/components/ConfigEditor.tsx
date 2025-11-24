import { Save } from "lucide-react";
import type React from "react";
import { useCallback, useEffect, useState } from "react";
import type { config } from "../../wailsjs/go/config/models";
import { useConfig } from "../context/ConfigContext";
import "./ConfigEditor.css";

interface ConfigEditorProps {
  configItem: config.ConfigItem;
}

const ConfigEditor: React.FC<ConfigEditorProps> = ({ configItem }) => {
  const { readConfig, saveConfig } = useConfig();
  const [content, setContent] = useState<string>("");
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [isDirty, setIsDirty] = useState<boolean>(false);
  const [isSaving, setIsSaving] = useState<boolean>(false);

  const loadFile = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      // In mock mode (browser), this might fail if ReadConfig isn't mocked in window.go
      // Ideally we should handle that gracefully or use mock data
      const text = await readConfig(configItem.path);
      setContent(text);
      setIsDirty(false);
    } catch (err) {
      console.error("Error loading file:", err);
      setError(
        "Failed to load file content. Please check if the backend is running.",
      );
      // For demo purposes in browser environment without backend:
      if (String(err).includes("window.go")) {
        setContent(
          `// Mock content for ${configItem.name}\n// Backend not connected.`,
        );
        setIsDirty(false);
        setError(null);
      }
    } finally {
      setIsLoading(false);
    }
  }, [configItem.path, configItem.name, readConfig]);

  useEffect(() => {
    loadFile();
  }, [loadFile]);

  const handleSave = async () => {
    setIsSaving(true);
    try {
      // Basic JSON validation
      if (configItem.format === "json") {
        try {
          JSON.parse(content);
        } catch (_e) {
          alert("Invalid JSON format. Please fix errors before saving.");
          setIsSaving(false);
          return;
        }
      }

      await saveConfig(configItem.path, content);
      setIsDirty(false);
    } catch (err) {
      console.error("Error saving file:", err);
      alert(`Failed to save file: ${err}`);
    } finally {
      setIsSaving(false);
    }
  };

  const handleTextChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setContent(e.target.value);
    setIsDirty(true);
  };

  return (
    <div className="config-editor">
      <div className="editor-toolbar">
        <div className="file-info">
          <span className="file-name">{configItem.name}</span>
          <span className="file-path">{configItem.path}</span>
        </div>
        <div className="editor-actions">
          <button
            type="button"
            className="btn-save"
            onClick={handleSave}
            disabled={!isDirty || isSaving || isLoading}
          >
            <Save size={16} />
            {isSaving ? "Saving..." : "Save"}
          </button>
        </div>
      </div>

      <div className="editor-area">
        {isLoading ? (
          <div className="editor-loading">Loading...</div>
        ) : error ? (
          <div className="editor-error">{error}</div>
        ) : (
          <textarea
            className="editor-textarea"
            value={content}
            onChange={handleTextChange}
            spellCheck={false}
          />
        )}
      </div>
    </div>
  );
};

export default ConfigEditor;
