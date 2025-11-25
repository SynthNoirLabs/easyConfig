import { useState, useEffect } from "react";
import { toast } from "sonner";
import { FetchPopularServers, InstallMCPPackage } from "../../wailsjs/go/main/App";
import { Search, Download, Star, Server, ExternalLink } from "lucide-react";
import "./Marketplace.css";

interface MCPPackage {
  name: string;
  description: string;
  stars: number;
  downloads: number;
  author: string;
  tags: string[];
}

export default function Marketplace() {
  const [servers, setServers] = useState<MCPPackage[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchQuery, setSearchQuery] = useState("");
  const [installing, setInstalling] = useState<string | null>(null);

  useEffect(() => {
    loadServers();
  }, []);

  const loadServers = async () => {
    try {
      const data = await FetchPopularServers();
      // Map backend data to frontend interface if needed, or use directly
      // Assuming backend returns compatible structure
      setServers(data || []);
    } catch (err) {
      toast.error("Failed to load marketplace servers");
    } finally {
      setLoading(false);
    }
  };

  const handleInstall = async (pkg: MCPPackage) => {
    setInstalling(pkg.name);
    try {
      await InstallMCPPackage(pkg.name);
      toast.success(`Successfully installed ${pkg.name}`);
    } catch (err) {
      toast.error(`Failed to install ${pkg.name}: ${err}`);
    } finally {
      setInstalling(null);
    }
  };

  const filteredServers = servers.filter(s => 
    s.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    s.description.toLowerCase().includes(searchQuery.toLowerCase())
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
        </div>
      </div>

      {loading ? (
        <div className="loading-state">
          <div className="spinner"></div>
          <p>Loading servers...</p>
        </div>
      ) : (
        <div className="servers-grid">
          {filteredServers.map((server, i) => (
            <div key={i} className="server-card">
              <div className="server-header">
                <div className="server-icon">
                  <Server size={24} />
                </div>
                <div className="server-meta">
                  <h3>{server.name}</h3>
                  <span className="server-author">by {server.author || "Unknown"}</span>
                </div>
              </div>
              
              <p className="server-desc">{server.description}</p>
              
              <div className="server-stats">
                {server.stars > 0 && (
                  <div className="stat" title="Stars">
                    <Star size={14} />
                    <span>{server.stars}</span>
                  </div>
                )}
                {server.downloads > 0 && (
                  <div className="stat" title="Downloads">
                    <Download size={14} />
                    <span>{server.downloads}</span>
                  </div>
                )}
              </div>

              <div className="server-tags">
                {server.tags && server.tags.slice(0, 3).map(tag => (
                  <span key={tag} className="tag">{tag}</span>
                ))}
              </div>

              <button 
                className="btn-install"
                onClick={() => handleInstall(server)}
                disabled={installing === server.name}
              >
                {installing === server.name ? "Installing..." : "Install"}
              </button>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
