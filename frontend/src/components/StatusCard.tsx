import {
  AlertTriangle,
  CheckCircle,
  FileText,
  HelpCircle,
  XCircle,
} from "lucide-react";
import type React from "react";
import { useState } from "react";
import type { config } from "../../wailsjs/go/models";
import ProviderStatusDrawer from "./ProviderStatusDrawer";
import "./StatusCard.css";

interface StatusCardProps {
  status: config.ProviderStatus;
}

const StatusCard: React.FC<StatusCardProps> = ({ status }) => {
  const [isDrawerOpen, setIsDrawerOpen] = useState(false);

  const getStatusIcon = (health: string) => {
    switch (health) {
      case "healthy":
        return <CheckCircle className="status-icon healthy" size={20} />;
      case "warning":
        return <AlertTriangle className="status-icon warning" size={20} />;
      case "error":
        return <XCircle className="status-icon error" size={20} />;
      default:
        return <HelpCircle className="status-icon unknown" size={20} />;
    }
  };

  const getStatusClass = (health: string) => {
    return `status-card ${health}`;
  };

  return (
    <>
      <div
        className={getStatusClass(status.health)}
        onClick={() => setIsDrawerOpen(true)}
        role="button"
        tabIndex={0}
        onKeyDown={(e) => {
          if (e.key === "Enter" || e.key === " ") {
            setIsDrawerOpen(true);
          }
        }}
      >
        <div className="status-card-header">
          {getStatusIcon(status.health)}
          <h3>{status.providerName}</h3>
        </div>
        <div className="status-message">
          <p>{status.statusMessage || "Operational"}</p>
        </div>
        <div className="status-meta">
          <div className="meta-item">
            <FileText size={14} />
            <span>{status.discoveredFiles?.length || 0} Files</span>
          </div>
          <div className="meta-item">
            <span>Last checked: {new Date(status.lastChecked).toLocaleTimeString()}</span>
          </div>
        </div>
      </div>
      <ProviderStatusDrawer
        isOpen={isDrawerOpen}
        onClose={() => setIsDrawerOpen(false)}
        status={status}
      />
    </>
  );
};

export default StatusCard;
