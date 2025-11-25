import type React from "react";
import { useEffect, useState } from "react";
import { toast } from "sonner";
import { FetchPopularServers, InstallMCPPackage } from "../../wailsjs/go/main/App";
import { marketplaces } from "../../wailsjs/go/models";
import "./Marketplace.css";

const Marketplace: React.FC = () => {
  const [packages, setPackages] = useState<marketplaces.MCPPackage[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [installing, setInstalling] = useState<string | null>(null);

  useEffect(() => {
    const loadPackages = async () => {
      try {
        const data = await FetchPopularServers();
        setPackages(data || []);
      } catch (err) {
        console.error("Failed to fetch packages:", err);
        setError("Failed to load marketplace data.");
        toast.error("Failed to load marketplace data.");
      } finally {
        setLoading(false);
      }
    };

    loadPackages();
  }, []);

  const handleInstall = async (pkg: marketplaces.MCPPackage) => {
    setInstalling(pkg.name);
    toast.info(`Installing ${pkg.name}...`);
    try {
      await InstallMCPPackage(pkg);
      toast.success(`Successfully installed ${pkg.name}`);
    } catch (err) {
      console.error("Failed to install package:", err);
      toast.error(`Failed to install ${pkg.name}`);
    } finally {
      setInstalling(null);
    }
  };

  if (loading) {
    return <div className="marketplace-loading">Loading marketplace...</div>;
  }

  if (error) {
    return <div className="marketplace-error">{error}</div>;
  }

  return (
    <div className="marketplace">
      <h2>MCP Marketplace</h2>
      <div className="marketplace-grid">
        {packages.map((pkg) => (
          <div key={pkg.name} className="marketplace-card">
            <h3>{pkg.name}</h3>
            <p className="marketplace-desc">{pkg.description}</p>
            <div className="marketplace-meta">
              <span>v{pkg.version}</span>
              <span>by {pkg.author}</span>
            </div>
            <button
              className="btn-install"
              onClick={() => handleInstall(pkg)}
              disabled={installing === pkg.name}
            >
              {installing === pkg.name ? "Installing..." : "Install"}
            </button>
          </div>
        ))}
      </div>
    </div>
  );
};

export default Marketplace;
