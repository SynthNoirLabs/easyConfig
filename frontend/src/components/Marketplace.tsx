import {
  Download,
  Filter,
  Loader2,
  RefreshCw,
  Search,
  Server,
  Star,
  Tag,
  Verified,
} from "lucide-react";
import type React from "react";
import { useCallback, useEffect, useState } from "react";
import { toast } from "sonner";
import {
  FetchPopularServers,
  InstallMCPPackage,
} from "../../wailsjs/go/main/App";
import type { marketplaces } from "../../wailsjs/go/models";
import "./Marketplace.css";

const Marketplace: React.FC = () => {
  const [packages, setPackages] = useState<marketplaces.MCPPackage[]>([]);
  const [filteredPackages, setFilteredPackages] = useState<
    marketplaces.MCPPackage[]
  >([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [installing, setInstalling] = useState<string | null>(null);
  const [searchTerm, setSearchTerm] = useState<string>("");
  const [selectedTag, setSelectedTag] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [cacheStatus, setCacheStatus] = useState<{
    isCached: boolean;
    isStale: boolean;
  }>({ isCached: false, isStale: true });
  const [refreshingCache, setRefreshingCache] = useState(false);

  // Fetch marketplace data
  const fetchMarketplace = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      // Mock cache status check since backend method is missing
      // const status = await GetMarketplaceCacheStatus();
      // setCacheStatus(status);
      setCacheStatus({ isCached: true, isStale: false });

      const results = await FetchPopularServers();
      setPackages(results || []);
      setFilteredPackages(results || []);
    } catch (err) {
      console.error("Failed to fetch marketplace data:", err);
      setError("Failed to load marketplace packages. Please try again later.");
    } finally {
      setLoading(false);
    }
  }, []);

  // Initial load
  useEffect(() => {
    fetchMarketplace();
  }, [fetchMarketplace]);

  // Handle refresh cache
  const handleRefreshCache = async () => {
    setRefreshingCache(true);
    try {
      // Mock refresh since backend method is missing
      // await RefreshMarketplaceCache();
      await fetchMarketplace();
      toast.success("Marketplace cache refreshed");
    } catch (err) {
      toast.error("Failed to refresh cache");
    } finally {
      setRefreshingCache(false);
    }
  };

  // Filter logic
  useEffect(() => {
    let result = packages;

    if (searchTerm) {
      const lowerTerm = searchTerm.toLowerCase();
      result = result.filter(
        (pkg) =>
          pkg.name.toLowerCase().includes(lowerTerm) ||
          pkg.description.toLowerCase().includes(lowerTerm) ||
          pkg.tags?.some((tag) => tag.toLowerCase().includes(lowerTerm)),
      );
    }

    if (selectedTag) {
      result = result.filter((pkg) => pkg.tags?.includes(selectedTag));
    }

    setFilteredPackages(result);
  }, [searchTerm, selectedTag, packages]);

  // Handle install
  const handleInstall = async (pkg: marketplaces.MCPPackage) => {
    setInstalling(pkg.name);
    try {
      await InstallMCPPackage(pkg.source); // Assuming source is the install arg
      toast.success(`Successfully installed ${pkg.name}`);
    } catch (err) {
      console.error(`Failed to install ${pkg.name}:`, err);
      toast.error(`Failed to install ${pkg.name}`);
    } finally {
      setInstalling(null);
    }
  };

  // Extract all unique tags
  const allTags = Array.from(
    new Set(packages.flatMap((pkg) => pkg.tags || [])),
  ).sort();

  if (error) {
    return (
      <div className="marketplace-error">
        <div className="error-content">
          <Server size={48} />
          <h3>Connection Error</h3>
          <p>{error}</p>
          <button type="button" onClick={fetchMarketplace} className="btn-retry">
            Try Again
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="marketplace-container">
      <div className="marketplace-header">
        <div className="header-title">
          <h1>MCP Marketplace</h1>
          <p>Discover and install extensions for your AI agents</p>
        </div>
        <div className="header-actions">
          {cacheStatus.isStale && (
            <div className="cache-warning" title="Data might be outdated">
              <RefreshCw size={14} /> Stale Cache
            </div>
          )}
          <button
            type="button"
            className="btn-refresh"
            onClick={handleRefreshCache}
            disabled={refreshingCache || loading}
            title="Refresh Marketplace Data"
          >
            <RefreshCw
              size={18}
              className={refreshingCache ? "spin" : ""}
            />
            Refresh
          </button>
        </div>
      </div>

      <div className="marketplace-filters">
        <div className="search-box">
          <Search size={18} />
          <input
            type="text"
            placeholder="Search packages..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
          />
        </div>

        <div className="filter-tags">
          <Filter size={16} className="filter-icon" />
          <button
            type="button"
            className={`tag-chip ${selectedTag === null ? "active" : ""}`}
            onClick={() => setSelectedTag(null)}
          >
            All
          </button>
          {allTags.map((tag) => (
            <button
              key={tag}
              type="button"
              className={`tag-chip ${selectedTag === tag ? "active" : ""}`}
              onClick={() => setSelectedTag(tag)}
            >
              {tag}
            </button>
          ))}
        </div>
      </div>

      {loading ? (
        <div className="marketplace-loading">
          <Loader2 size={40} className="spin" />
          <p>Loading marketplace data...</p>
        </div>
      ) : (
        <div className="packages-grid">
          {filteredPackages.length > 0 ? (
            filteredPackages.map((pkg) => (
              <div key={pkg.name} className="package-card">
                <div className="package-header">
                  <div className="package-icon">
                    {pkg.vendor ? (
                      <img
                        src={`https://github.com/${pkg.vendor}.png`}
                        alt={pkg.vendor}
                        onError={(e) => {
                          (e.target as HTMLImageElement).style.display = "none";
                        }}
                      />
                    ) : (
                      <Server size={24} />
                    )}
                  </div>
                  <div className="package-title-row">
                    <h3>{pkg.name}</h3>
                    {pkg.verified && (
                      <Verified
                        size={16}
                        className="verified-badge"

                      />
                    )}
                  </div>
                </div>

                <p className="package-description">{pkg.description}</p>

                <div className="package-meta">
                  {pkg.stars !== undefined && (
                    <div className="meta-stat">
                      <Star size={14} /> {pkg.stars}
                    </div>
                  )}
                  {pkg.downloads !== undefined && (
                    <div className="meta-stat">
                      <Download size={14} /> {pkg.downloads}
                    </div>
                  )}
                  <div className="meta-stat">v{pkg.version || "latest"}</div>
                </div>

                <div className="package-tags">
                  {pkg.tags?.slice(0, 3).map((tag) => (
                    <span key={tag} className="pkg-tag">
                      <Tag size={10} /> {tag}
                    </span>
                  ))}
                  {(pkg.tags?.length || 0) > 3 && (
                    <span className="pkg-tag more">
                      +{ (pkg.tags?.length || 0) - 3 }
                    </span>
                  )}
                </div>

                <div className="package-actions">
                  <button
                    type="button"
                    className="btn-install"
                    onClick={() => handleInstall(pkg)}
                    disabled={installing === pkg.name}
                  >
                    {installing === pkg.name ? (
                      <>
                        <Loader2 size={16} className="spin" /> Installing...
                      </>
                    ) : (
                      <>
                        <Download size={16} /> Install
                      </>
                    )}
                  </button>
                  <a
                    href={pkg.repoUrl || pkg.url}
                    target="_blank"
                    rel="noreferrer"
                    className="btn-view"
                  >
                    View
                  </a>
                </div>
              </div>
            ))
          ) : (
            <div className="no-results">
              <Search size={48} />
              <p>No packages found matching your criteria</p>
              <button
                type="button"
                className="btn-secondary"
                onClick={() => {
                  setSearchTerm("");
                  setSelectedTag(null);
                }}
              >
                Clear Filters
              </button>
            </div>
          )}
        </div>
      )}
    </div>
  );
};

export default Marketplace;
