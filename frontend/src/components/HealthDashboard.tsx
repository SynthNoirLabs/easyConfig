import { RefreshCw } from "lucide-react";
import type React from "react";
import { useCallback, useEffect, useState } from "react";
import { GetProviderStatuses } from "../../wailsjs/go/main/App";
import type { config } from "../../wailsjs/go/models";
import "./HealthDashboard.css";
import StatusCard from "./StatusCard";

const HealthDashboard: React.FC = () => {
  const [statuses, setStatuses] = useState<config.ProviderStatus[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  const fetchStatuses = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const result = await GetProviderStatuses();
      setStatuses(result);
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Failed to fetch status reports",
      );
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchStatuses();
    const interval = setInterval(fetchStatuses, 30000); // Auto-refresh every 30s
    return () => clearInterval(interval);
  }, [fetchStatuses]);

  if (error) {
    return (
      <div className="health-dashboard-error">
        <p>Error loading dashboard: {error}</p>
        <button type="button" onClick={fetchStatuses}>Retry</button>
      </div>
    );
  }

  return (
    <div className="health-dashboard">
      <div className="dashboard-header">
        <h1>Provider Health Dashboard</h1>
        <button type="button" onClick={fetchStatuses} disabled={loading}>
          <RefreshCw size={16} />
          {loading ? "Refreshing..." : "Refresh All"}
        </button>
      </div>

      <div className="status-grid">
        {statuses.map((status) => (
          <StatusCard key={status.providerName} status={status} />
        ))}
      </div>
    </div>
  );
};

export default HealthDashboard;
