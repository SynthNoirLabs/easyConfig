import Editor from "@monaco-editor/react";
import { Save } from "lucide-react";
import type React from "react";
import { useCallback, useEffect, useState } from "react";
import { toast } from "sonner"; // Import sonner toast
import type { config } from "../../wailsjs/go/config/models";
import { useConfig } from "../context/ConfigContext";
import "./ConfigEditor.css";

interface ConfigEditorProps {
  configItem: config.ConfigItem;
}

const getLanguage = (format: string) => {
  switch (format.toLowerCase()) {
    case "json":
      return "json";
    case "yaml":
    case "yml":
      return "yaml";
    case "toml":
      return "ini"; // TOML syntax highlighting is not built-in, INI is closest
    case "ini":
      return "ini";
    default:
      return "plaintext";
  }
};

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
      const text = await readConfig(configItem.path);
      setContent(text);
      setIsDirty(false);
    } catch (err) {
      console.error("Error loading file:", err);
      toast.error("Failed to load file content."); // Use toast for error
      setError(
        err instanceof Error ? err.message : "Failed to load configurations",
      );
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
      if (configItem.format === "json") {
        try {
          JSON.parse(content);
        } catch (_e) {
          toast.error("Invalid JSON format. Please fix errors before saving."); // Use toast for error
          setIsSaving(false);
          return;
        }
      }

      await saveConfig(configItem.path, content);
      setIsDirty(false);
      toast.success("Configuration saved successfully!"); // Use toast for success
    } catch (err) {
      console.error("Error saving file:", err);
      toast.error(err instanceof Error ? err.message : "Failed to save file."); // Use toast for error
    } finally {
      setIsSaving(false);
    }
  };

  const handleEditorChange = (value: string | undefined) => {
    setContent(value || "");
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
          <Editor
            height="100%"
            defaultLanguage="plaintext"
            language={getLanguage(configItem.format)}
            value={content}
            theme="vs-dark"
            onChange={handleEditorChange}
            options={{
              minimap: { enabled: false },
              scrollBeyondLastLine: false,
              fontSize: 14,
              automaticLayout: true,
            }}
          />
        )}
      </div>
    </div>
  );
};

export default ConfigEditor;
