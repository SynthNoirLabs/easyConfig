import {
  BadgeCheck,
  Download,
  Link,
  Search,
  SearchX,
  Server,
  ShieldCheck,
  Star,
} from "lucide-react";
import { useCallback, useEffect, useMemo, useState } from "react";
import { toast } from "sonner";
import {
  FetchPopularServers,
  InstallMCPPackage,
} from "../../wailsjs/go/main/App";
import type { marketplaces } from "../../wailsjs/go/models";
import { EmptyState } from "./EmptyState";
import "./Marketplace.css";

export default function Marketplace() {
  const [servers, setServers] = useState<marketplaces.MCPPackage[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchQuery, setSearchQuery] = useState("");
  const [installing, setInstalling] = useState<string | null>(null);
  const [verifiedOnly, setVerifiedOnly] = useState(false);

  const loadServers = useCallback(async () => {
    try {
      const data = await FetchPopularServers();
      setServers(data || []);
    } catch (_err) {
      toast.error("Failed to load marketplace servers");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    loadServers();
  }, [loadServers]);

  const handleInstall = async (pkg: marketplaces.MCPPackage) => {
    setInstalling(pkg.name);
    try {
      await InstallMCPPackage(JSON.stringify(pkg));
      toast.success(`Successfully installed ${pkg.name}`);
    } catch (err) {
      toast.error(`Failed to install ${pkg.name}: ${err}`);
    } finally {
      setInstalling(null);
    }
  };

  const filteredServers = servers.filter((s) => {
    if (verifiedOnly && !s.verified) return false;
    const q = searchQuery.toLowerCase();
    return (
      s.name.toLowerCase().includes(q) ||
      s.description.toLowerCase().includes(q) ||
      (s.tags || []).some((t) => t.toLowerCase().includes(q))
    );
  });

  const sortedServers = useMemo(
    () =>
      [...filteredServers].sort(
        (a, b) => Number(b.verified) - Number(a.verified),
      ),
    [filteredServers],
  );

  return (
    <div className="marketplace-container">
      <div className="marketplace-header">
        <div>
          <h2>MCP Marketplace</h2>
          <p>Discover and install Model Context Protocol servers.</p>
        </div>
        <div className="search-box">
          <Search size={18} className="search-icon" />
          <input
            type="text"
            placeholder="Search servers..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="search-input"
          />
          <button
            type="button"
            className={`filter-chip ${verifiedOnly ? "active" : ""}`}
            onClick={() => setVerifiedOnly((v) => !v)}
            title="Show only verified packages"
          >
            <ShieldCheck size={14} /> Verified only
          </button>
        </div>
      </div>

      {loading ? (
        <div className="loading-state">
          <div className="spinner"></div>
          <p>Loading servers...</p>
        </div>
      ) : sortedServers.length > 0 ? (
        <div className="servers-grid">
          {sortedServers.map((server) => (
            <div
              key={`${server.source || "smithery"}-${server.name}`}
              className="server-card"
            >
              <div className="server-header">
                <div className="server-icon">
                  <Server size={24} />
                </div>
                <div className="server-meta">
                  <h3>{server.name}</h3>
                  <span className="server-author">
                    by {server.author || "Unknown"}
                  </span>
                  <div className="server-meta-row">
                    {server.license && (
                      <span className="meta-pill">{server.license}</span>
                    )}
                    {server.verified && (
                      <span className="meta-pill verified">
                        <BadgeCheck size={14} /> Verified
                      </span>
                    )}
                  </div>
                </div>
              </div>

              <p className="server-desc">{server.description}</p>

              <div className="server-stats">
                {(server.stars || 0) > 0 && (
                  <div className="stat" title="Stars">
                    <Star size={14} />
                    <span>{server.stars}</span>
                  </div>
                )}
                {(server.downloads || 0) > 0 && (
                  <div className="stat" title="Downloads">
                    <Download size={14} />
                    <span>{server.downloads}</span>
                  </div>
                )}
              </div>

              <div className="server-tags">
                {server.tags?.slice(0, 3).map((tag) => (
                  <span key={tag} className="tag">
                    {tag}
                  </span>
                ))}
              </div>

              <div className="server-links">
                {server.repoUrl && (
                  <a
                    className="link"
                    href={server.repoUrl}
                    target="_blank"
                    rel="noreferrer"
                  >
                    <Link size={14} /> Repo
                  </a>
                )}
                {server.url && (
                  <a
                    className="link"
                    href={server.url}
                    target="_blank"
                    rel="noreferrer"
                  >
                    <Server size={14} /> Homepage
                  </a>
                )}
              </div>

              <button
                type="button"
                className="btn-install"
                onClick={() => handleInstall(server)}
                disabled={installing === server.name}
              >
                {installing === server.name ? "Installing..." : "Install"}
              </button>
            </div>
          ))}
        </div>
      ) : (
        <EmptyState
          icon={SearchX}
          title="No Results"
          description="Try adjusting your search or filters."
        />
      )}
    </div>
  );
}
