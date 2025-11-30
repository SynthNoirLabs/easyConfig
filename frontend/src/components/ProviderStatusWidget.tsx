import type React from "react";
import { useEffect, useState } from "react";
import { GetProviderStatuses } from "../../wailsjs/go/main/App";
import ProviderStatusDrawer from "./ProviderStatusDrawer";
import "./ProviderStatusWidget.css";

// Define the type based on the Go struct
interface ProviderStatus {
  providerName: string;
  health: "healthy" | "unhealthy" | "unknown";
  statusMessage: string;
  lastChecked: string;
}

const ProviderStatusWidget: React.FC = () => {
  const [statuses, setStatuses] = useState<ProviderStatus[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [selectedStatus, setSelectedStatus] = useState<ProviderStatus | null>(
    null,
  );

  useEffect(() => {
    const fetchStatuses = async () => {
      try {
        const result = await GetProviderStatuses();
        setStatuses(result);
      } catch (err) {
        setError("Failed to fetch provider statuses.");
        console.error(err);
      }
    };

    fetchStatuses();
    // Refresh statuses every 30 seconds
    const interval = setInterval(fetchStatuses, 30000);

    return () => clearInterval(interval);
  }, []);

  const handlePillClick = (status: ProviderStatus) => {
    setSelectedStatus(status);
  };

  const handleDrawerClose = () => {
    setSelectedStatus(null);
  };

  const getHealthColor = (health: "healthy" | "unhealthy" | "unknown") => {
    switch (health) {
      case "healthy":
        return "var(--green)";
      case "unhealthy":
        return "var(--red)";
      default:
        return "var(--gray)";
    }
  };

  if (error) {
    return <div className="provider-status-widget error">{error}</div>;
  }

  return (
    <>
      <div className="provider-status-widget">
        <h4>Provider Status</h4>
        <div className="status-pills">
          {statuses.map((status) => (
            <button
              type="button"
              key={status.providerName}
              className="status-pill"
              style={{ backgroundColor: getHealthColor(status.health) }}
              title={`${status.providerName}: ${status.statusMessage}`}
              onClick={() => handlePillClick(status)}
              onKeyDown={(e) => {
                if (e.key === "Enter" || e.key === " ") {
                  handlePillClick(status);
                }
              }}
            >
              {status.providerName}
            </button>
          ))}
        </div>
      </div>
      <ProviderStatusDrawer
        status={selectedStatus}
        onClose={handleDrawerClose}
      />
    </>
  );
};

export default ProviderStatusWidget;
