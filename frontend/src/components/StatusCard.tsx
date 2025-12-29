import { AlertTriangle, CheckCircle, HelpCircle, XCircle } from "lucide-react";
import type React from "react";
import "./StatusCard.css";
import type { config } from "../../wailsjs/go/models";

interface StatusCardProps {
  status: config.ProviderStatusReport;
  onClick: () => void;
}

const getStatusLevel = (status: config.ProviderStatusReport) => {
  if (status.installed && status.configured) return "ready";
  if (status.installed && !status.configured) return "warning";
  if (!status.installed) return "not-installed";
  return "error";
};

const StatusIcon = ({ status }: { status: config.ProviderStatusReport }) => {
  const level = getStatusLevel(status);
  switch (level) {
    case "ready":
      return <CheckCircle className="status-icon ready" />;
    case "warning":
      return <AlertTriangle className="status-icon warning" />;
    case "error":
      return <XCircle className="status-icon error" />;
    case "not-installed":
      return <HelpCircle className="status-icon not-installed" />;
    default:
      return <HelpCircle className="status-icon not-installed" />;
  }
};

const StatusCard: React.FC<StatusCardProps> = ({ status, onClick }) => {
  const level = getStatusLevel(status);

  return (
    <div className={`status-card ${level}`} onClick={onClick}>
      <StatusIcon status={status} />
      <h4>{status.providerName}</h4>
      <p>{status.message}</p>
      {status.version && <span className="version">v{status.version}</span>}
    </div>
  );
};

export default StatusCard;
