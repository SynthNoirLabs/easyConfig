import { X } from "lucide-react";
import type React from "react";
import type { config } from "../../wailsjs/go/models";
import "./ProviderStatusDrawer.css";

interface ProviderStatusDrawerProps {
  isOpen: boolean;
  onClose: () => void;
  status: config.ProviderStatus;
}

const ProviderStatusDrawer: React.FC<ProviderStatusDrawerProps> = ({
  isOpen,
  onClose,
  status,
}) => {
  if (!isOpen) return null;

  return (
    <div className="status-drawer-overlay" onClick={onClose}
      onKeyDown={(e) => {
        if (e.key === "Escape") onClose();
      }}
      role="button"
      tabIndex={0}
    >
      <div
        className="status-drawer"
        onClick={(e) => e.stopPropagation()}
        role="dialog"
        aria-modal="true"
      >
        <div className="status-drawer-header">
          <h3>{status.providerName} Status</h3>
          <button type="button" className="close-btn" onClick={onClose}>
            <X size={20} />
          </button>
        </div>
        <div className="status-drawer-content">
          <div className="status-summary">
            <p>
              <strong>Health:</strong>{" "}
              <span className={`status-badge ${status.health}`}>
                {status.health.toUpperCase()}
              </span>
            </p>
            <p>
              <strong>Last Checked:</strong>{" "}
              {new Date(status.lastChecked).toLocaleString()}
            </p>
            <p>
              <strong>Message:</strong>{" "}
              {status.statusMessage || "All systems operational."}
            </p>
          </div>

          <div className="discovered-files">
            <h4>Discovered Files ({status.discoveredFiles?.length || 0})</h4>
            {status.discoveredFiles && status.discoveredFiles.length > 0 ? (
              <ul>
                {status.discoveredFiles.map((file) => (
                  <li key={file.path}>
                    <div className="file-info">
                      <span className="file-name">{file.name}</span>
                      <span className="file-path">{file.path}</span>
                    </div>
                  </li>
                ))}
              </ul>
            ) : (
              <p className="no-files">No configuration files found.</p>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default ProviderStatusDrawer;
