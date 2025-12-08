import type React from "react";
import { useEffect, useState } from "react";
import { GetProviderStatuses } from "../../wailsjs/go/main/App";
import ProviderStatusDrawer from "./ProviderStatusDrawer";
import ConfigWizard from "./ConfigWizard";
import type { config } from "../../wailsjs/go/models";
import "./ProviderStatusWidget.css";

type ProviderStatus = config.ProviderStatus;

const ProviderStatusWidget: React.FC = () => {
  const [statuses, setStatuses] = useState<ProviderStatus[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [selectedStatus, setSelectedStatus] = useState<ProviderStatus | null>(
    null,
  );
  const [wizardProvider, setWizardProvider] = useState<string | null>(null);

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

  const getHealthColor = (health: string) => {
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
            <div key={status.providerName} className="status-pill-container">
              <button
                type="button"
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
              {status.hasWizard && (
                <button
                  type="button"
                  className="wizard-button"
                  onClick={() => setWizardProvider(status.providerName)}
                >
                  Wizard
                </button>
              )}
            </div>
          ))}
        </div>
      </div>
      <ProviderStatusDrawer
        status={selectedStatus}
        onClose={handleDrawerClose}
      />
      {wizardProvider && (
        <ConfigWizard
          providerName={wizardProvider}
          isOpen={!!wizardProvider}
          onClose={() => setWizardProvider(null)}
        />
      )}
    </>
  );
};

export default ProviderStatusWidget;
