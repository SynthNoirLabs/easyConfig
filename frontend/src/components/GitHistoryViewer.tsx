import type React from "react";
import { versions } from "../../wailsjs/go/models";
import "./GitHistoryViewer.css";

interface GitHistoryViewerProps {
  history: versions.CommitInfo[];
  onSelectCommit: (commitHash: string) => void;
  onClose: () => void;
}

const GitHistoryViewer: React.FC<GitHistoryViewerProps> = ({
  history,
  onSelectCommit,
  onClose,
}) => {
  return (
    <div className="git-history-viewer-overlay">
      <div className="git-history-viewer">
        <div className="history-header">
          <h2>Git History</h2>
          <button type="button" className="btn-close" onClick={onClose}>
            &times;
          </button>
        </div>
        <div className="history-list">
          {history.length === 0 ? (
            <p>No git history found for this file.</p>
          ) : (
            <ul>
              {history.map((commit) => (
                <li key={commit.hash}>
                  <div className="commit-info">
                    <p className="commit-message">{commit.message}</p>
                    <p className="commit-meta">
                      <strong>{commit.author}</strong> on{" "}
                      {new Date(commit.date).toLocaleString()}
                    </p>
                    <p className="commit-hash">{commit.hash.substring(0, 7)}</p>
                  </div>
                  <div className="commit-actions">
                    <button
                      type="button"
                      className="btn-revert"
                      onClick={() => onSelectCommit(commit.hash)}
                    >
                      Revert to this version
                    </button>
                  </div>
                </li>
              ))}
            </ul>
          )}
        </div>
      </div>
    </div>
  );
};

export default GitHistoryViewer;
