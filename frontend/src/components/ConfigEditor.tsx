import type { editor } from "monaco-editor";
import Editor from "@monaco-editor/react";
import {
  Code,
  Eye,
  GitCompare,
  LayoutTemplate,
  RefreshCw,
  RotateCcw,
  Save,
} from "lucide-react";
import type React from "react";
import { useCallback, useEffect, useRef, useState } from "react";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import { toast } from "sonner"; // Import sonner toast
import type { config } from "../../wailsjs/go/models";
import { useConfig } from "../context/ConfigContext";
import { useKeyboardShortcuts } from "../hooks/useKeyboardShortcuts";
import "./ConfigEditor.css";
import DiffViewer from "./DiffViewer";
import ClaudeConfigEditor from "./editors/ClaudeConfigEditor";
import OpenCodeConfigEditor from "./editors/OpenCodeConfigEditor";

interface ConfigEditorProps {
  configItem: config.Item & { initialLine?: number };
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
  const [originalContent, setOriginalContent] = useState<string>("");
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [isDirty, setIsDirty] = useState<boolean>(false);
  const [isSaving, setIsSaving] = useState<boolean>(false);
  const [viewMode, setViewMode] = useState<
    "code" | "form" | "preview" | "compare"
  >("code");
  const editorRef = useRef<editor.IStandaloneCodeEditor | null>(null);

  const isMarkdown =
    configItem.format.toLowerCase() === "markdown" ||
    configItem.fileName.toLowerCase().endsWith(".md");

  // Determine if a specific editor is available
  const hasSpecificEditor =
    (configItem.provider === "Claude Code" &&
      configItem.fileName === "claude_desktop_config.json") ||
    (configItem.provider === "OpenCode" &&
      configItem.fileName === "opencode.json");

  const handleSave = useCallback(async () => {
    if (!isDirty) return;
    setIsSaving(true);
    try {
      if (configItem.format === "json") {
        try {
          JSON.parse(content);
        } catch (_e) {
          toast.error("Invalid JSON format. Please fix errors before saving.");
          setIsSaving(false);
          return;
        }
      }

      await saveConfig(configItem.path, content);
      setOriginalContent(content);
      setIsDirty(false);
      toast.success("Configuration saved successfully!");
    } catch (err) {
      console.error("Error saving file:", err);
      toast.error(err instanceof Error ? err.message : "Failed to save file.");
    } finally {
      setIsSaving(false);
    }
  }, [content, configItem.path, saveConfig, isDirty]);

  useKeyboardShortcuts({
    "ctrl+s": handleSave,
    "ctrl+1": () => setViewMode("code"),
    "ctrl+2": () => hasSpecificEditor && setViewMode("form"),
    "ctrl+3": () => isMarkdown && setViewMode("preview"),
  });

  useEffect(() => {
    // Default to form view if available
    if (hasSpecificEditor) {
      setViewMode("form");
    } else if (isMarkdown) {
      setViewMode("preview");
    } else {
      setViewMode("code");
    }
  }, [hasSpecificEditor, isMarkdown]);

  const handleEditorDidMount = (editorInstance: editor.IStandaloneCodeEditor) => {
    editorRef.current = editorInstance;
    if (configItem.initialLine) {
      editorInstance.revealLineInCenter(configItem.initialLine);
      editorInstance.setPosition({
        lineNumber: configItem.initialLine,
        column: 1,
      });
    }
  };

  const loadFile = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const text = await readConfig(configItem.path);
      setContent(text);
      setOriginalContent(text);
      setIsDirty(false);
    } catch (err) {
      console.error("Error loading file:", err);
      toast.error("Failed to load file content.");
      setError(
        err instanceof Error ? err.message : "Failed to load configurations",
      );
      if (String(err).includes("window.go")) {
        const mock = `// Mock content for ${configItem.name}\n// Backend not connected.`;
        setContent(mock);
        setOriginalContent(mock);
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

  const handleReset = () => {
    if (confirm("Discard unsaved changes?")) {
      setContent(originalContent);
      setIsDirty(false);
    }
  };

  const handleReload = () => {
    if (
      isDirty &&
      !confirm(
        "You have unsaved changes. Reloading will discard them. Continue?",
      )
    ) {
      return;
    }
    loadFile();
  };

  const handleEditorChange = (value: string | undefined) => {
    setContent(value || "");
    setIsDirty(value !== originalContent);
  };

  return (
    <div className="config-editor">
      <div className="editor-toolbar">
        <div className="file-info">
          <span className="file-name">{configItem.name}</span>
          <span className="file-path">{configItem.path}</span>
        </div>
        <div className="editor-actions">
          {(hasSpecificEditor || isMarkdown || isDirty) && (
            <>
              <div className="view-toggle">
                <button
                  type="button"
                  className={`btn-toggle ${viewMode === "code" ? "active" : ""}`}
                  onClick={() => setViewMode("code")}
                  aria-label="Code View"
                >
                  <Code size={16} />
                </button>
                {hasSpecificEditor && (
                  <button
                    type="button"
                    className={`btn-toggle ${viewMode === "form" ? "active" : ""}`}
                    onClick={() => setViewMode("form")}
                    aria-label="Form View"
                  >
                    <LayoutTemplate size={16} />
                  </button>
                )}
                {isMarkdown && (
                  <button
                    type="button"
                    className={`btn-toggle ${viewMode === "preview" ? "active" : ""}`}
                    onClick={() => setViewMode("preview")}
                    aria-label="Preview"
                  >
                    <Eye size={16} />
                  </button>
                )}
                <button
                  type="button"
                  className={`btn-toggle ${viewMode === "compare" ? "active" : ""}`}
                  onClick={() => setViewMode("compare")}
                  title="Compare Changes"
                  disabled={!isDirty}
                >
                  <GitCompare size={16} />
                </button>
              </div>
              <div className="separator" />
            </>
          )}
          <button
            type="button"
            className="btn-secondary"
            onClick={handleReset}
            disabled={!isDirty || isLoading}
            aria-label="Reset to last saved"
          >
            <RotateCcw size={16} />
          </button>
          <button
            type="button"
            className="btn-secondary"
            onClick={handleReload}
            disabled={isLoading}
            aria-label="Reload from disk"
          >
            <RefreshCw size={16} />
          </button>
          <div className="separator" />
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
        ) : viewMode === "compare" ? (
          <DiffViewer
            original={originalContent}
            modified={content}
            language={getLanguage(configItem.format)}
          />
        ) : viewMode === "form" && hasSpecificEditor ? (
          configItem.provider === "Claude Code" ? (
            <ClaudeConfigEditor
              content={content}
              onChange={handleEditorChange}
            />
          ) : (
            <OpenCodeConfigEditor
              content={content}
              onChange={handleEditorChange}
            />
          )
        ) : viewMode === "preview" && isMarkdown ? (
          <div className="markdown-preview">
            <ReactMarkdown remarkPlugins={[remarkGfm]}>{content}</ReactMarkdown>
          </div>
        ) : (
          <Editor
            height="100%"
            defaultLanguage="plaintext"
            language={getLanguage(configItem.format)}
            value={content}
            theme="vs-dark"
            onChange={handleEditorChange}
            onMount={handleEditorDidMount}
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
