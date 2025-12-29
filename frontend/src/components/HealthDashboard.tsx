import { RefreshCw } from "lucide-react";
import type React from "react";
import { useCallback, useEffect, useState } from "react";
import { GetAllProviderStatuses } from "../../wailsjs/go/main/App";
import type { config } from "../../wailsjs/go/models";
import StatusCard from "./StatusCard";
import "./HealthDashboard.css";

const HealthDashboard: React.FC = () => {
  const [statuses, setStatuses] = useState<config.ProviderStatusReport[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchStatuses = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const result = await GetAllProviderStatuses();
      setStatuses(result);
    } catch (err) {
      setError("Failed to fetch provider statuses.");
      console.error(err);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchStatuses();
  }, [fetchStatuses]);

  const handleCardClick = (status: config.ProviderStatusReport) => {
    // For now, we'll just log the status.
    // In the future, this could open a detailed view.
    console.log("Card clicked:", status);
  };

  return (
    <div className="health-dashboard">
      <div className="dashboard-header">
        <h1>Provider Health Dashboard</h1>
        <button onClick={fetchStatuses} disabled={loading}>
          <RefreshCw size={16} />
          {loading ? "Refreshing..." : "Refresh All"}
        </button>
      </div>

      {loading && <p>Loading statuses...</p>}
      {error && <p className="error-message">{error}</p>}

      {!loading && !error && (
        <div className="status-grid">
          {statuses.map((status) => (
            <StatusCard
              key={status.providerName}
              status={status}
              onClick={() => handleCardClick(status)}
            />
          ))}
        </div>
      )}

      <div className="dashboard-legend">
        <span>✓ Ready</span>
        <span>⚠ Warning</span>
        <span>✗ Error</span>
        <span>○ Not Installed</span>
      </div>
    </div>
  );
};

export default HealthDashboard;
