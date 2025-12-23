import { X } from "lucide-react";
import type React from "react";
import type { config } from "../../wailsjs/go/models";
import "./ProviderStatusDrawer.css";

interface ProviderStatusDrawerProps {
  status: config.ProviderStatus | null;
  onClose: () => void;
}

const ProviderStatusDrawer: React.FC<ProviderStatusDrawerProps> = ({
  status,
  onClose,
}) => {
  if (!status) {
    return null;
  }

  const getOnboardingSteps = (providerName: string) => {
    // In a real app, this would be more dynamic
    switch (providerName.toLowerCase()) {
      case "claude code":
        return (
          <ul>
            <li>Create a global config file at ~/.claude/settings.json</li>
            <li>Add your API key to the configuration file.</li>
            <li>Refer to the Claude documentation for more details.</li>
          </ul>
        );
      case "gemini":
        return (
          <ul>
            <li>Create a global config file at ~/.gemini/settings.json</li>
            <li>Add your API key to the configuration file.</li>
          </ul>
        );
      default:
        return <p>No onboarding information available.</p>;
    }
  };

  return (
    <div className="status-drawer-overlay" onClick={onClose}>
      <div className="status-drawer" onClick={(e) => e.stopPropagation()}>
        <div className="status-drawer-header">
          <h3>{status.providerName} Status</h3>
          <button type="button" onClick={onClose} className="btn-icon">
            <X size={20} />
          </button>
        </div>
        <div className="status-drawer-content">
          <p>
            <strong>Health:</strong>{" "}
            <span className={`health-indicator ${status.health}`}>
              {status.health}
            </span>
          </p>
          <p>
            <strong>Message:</strong> {status.statusMessage}
          </p>
          <p>
            <strong>Last Checked:</strong>{" "}
            {new Date(status.lastChecked).toLocaleString()}
          </p>
          {status.health === "unhealthy" && (
            <div className="onboarding-checklist">
              <h4>Onboarding Checklist</h4>
              {getOnboardingSteps(status.providerName)}
            </div>
          )}
          <div className="quick-actions">
            <button type="button" className="btn">
              Open Config File
            </button>
            <button type="button" className="btn">
              Re-run Discovery
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ProviderStatusDrawer;
