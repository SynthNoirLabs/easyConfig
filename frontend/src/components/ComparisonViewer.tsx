import type React from "react";
import { useEffect, useState } from "react";
import { toast } from "sonner";
import type { config } from "../../wailsjs/go/models";
import { useConfig } from "../context/ConfigContext";
import DiffViewer from "./DiffViewer";

interface ComparisonViewerProps {
  item1: config.Item;
  item2: config.Item;
  onClose: () => void;
}

const getLanguage = (format: string) => {
  switch (format.toLowerCase()) {
    case "json":
      return "json";
    case "yaml":
    case "yml":
      return "yaml";
    case "toml":
      return "ini";
    case "ini":
      return "ini";
    default:
      return "plaintext";
  }
};

const ComparisonViewer: React.FC<ComparisonViewerProps> = ({
  item1,
  item2,
  onClose,
}) => {
  const { readConfig } = useConfig();
  const [content1, setContent1] = useState<string>("");
  const [content2, setContent2] = useState<string>("");
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchContents = async () => {
      setIsLoading(true);
      setError(null);
      try {
        const [text1, text2] = await Promise.all([
          readConfig(item1.path),
          readConfig(item2.path),
        ]);
        setContent1(text1);
        setContent2(text2);
      } catch (err) {
        console.error("Error loading files for comparison:", err);
        toast.error("Failed to load file contents for comparison.");
        setError(
          err instanceof Error
            ? err.message
            : "Failed to load file contents for comparison",
        );
      } finally {
        setIsLoading(false);
      }
    };

    fetchContents();
  }, [item1, item2, readConfig]);

  const language = getLanguage(item1.format);

  return (
    <div className="config-editor">
      {" "}
      {/* Reusing class for consistent styling */}
      <div className="editor-toolbar">
        <div className="file-info">
          <span className="file-name">
            Comparing: {item1.name} (
            <span className="file-path">{item1.path}</span>) vs {item2.name} (
            <span className="file-path">{item2.path}</span>)
          </span>
        </div>
        <div className="editor-actions">
          <button type="button" className="btn-secondary" onClick={onClose}>
            Close
          </button>
        </div>
      </div>
      <div className="editor-area">
        {isLoading ? (
          <div className="editor-loading">Loading comparison...</div>
        ) : error ? (
          <div className="editor-error">{error}</div>
        ) : (
          <DiffViewer
            original={content1}
            modified={content2}
            language={language}
          />
        )}
      </div>
    </div>
  );
};

export default ComparisonViewer;
