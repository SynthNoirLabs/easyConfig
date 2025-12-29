import type React from "react";
import { useEffect, useState } from "react";
import { toast } from "sonner";
import { ListProfileFiles, ReadConfig } from "../../wailsjs/go/main/App";
import DiffViewer from "./DiffViewer";

interface ProfilePreviewProps {
  profileName: string;
  onClose: () => void;
  onApply: () => void;
}

interface ProfileFileDiff {
  path: string;
  originalContent: string;
  modifiedContent: string;
  language: string;
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

const ProfilePreview: React.FC<ProfilePreviewProps> = ({
  profileName,
  onClose,
  onApply,
}) => {
  const [diffs, setDiffs] = useState<ProfileFileDiff[]>([]);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchProfileDiffs = async () => {
      setIsLoading(true);
      setError(null);
      try {
        const profileFiles = await ListProfileFiles(profileName);
        const diffPromises = profileFiles.map(async (file) => {
          const originalContent = await ReadConfig(file.path);
          return {
            path: file.path,
            originalContent,
            modifiedContent: file.content,
            language: getLanguage(file.path),
          };
        });
        const settledDiffs = await Promise.all(diffPromises);
        setDiffs(settledDiffs);
      } catch (err) {
        console.error("Error loading profile preview:", err);
        toast.error("Failed to load profile preview.");
        setError(
          err instanceof Error ? err.message : "Failed to load profile preview",
        );
      } finally {
        setIsLoading(false);
      }
    };
    fetchProfileDiffs();
  }, [profileName]);

  return (
    <div className="profile-preview-modal">
      <div className="profile-preview-header">
        <h2>Preview Profile: {profileName}</h2>
        <div className="profile-preview-actions">
          <button type="button" className="btn-secondary" onClick={onClose}>
            Cancel
          </button>
          <button type="button" className="btn-primary" onClick={onApply}>
            Apply Profile
          </button>
        </div>
      </div>
      <div className="profile-preview-content">
        {isLoading ? (
          <div className="editor-loading">Loading preview...</div>
        ) : error ? (
          <div className="editor-error">{error}</div>
        ) : (
          diffs.map((diff) => (
            <div key={diff.path} className="diff-container">
              <h3>{diff.path}</h3>
              <DiffViewer
                original={diff.originalContent}
                modified={diff.modifiedContent}
                language={diff.language}
              />
            </div>
          ))
        )}
      </div>
    </div>
  );
};

export default ProfilePreview;
