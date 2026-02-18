import {
  Activity,
  BookOpen,
  Brain,
  ChevronDown,
  ChevronRight,
  FileJson,
  Plus,
  Search,
  ShoppingBag,
  Terminal,
  Workflow,
} from "lucide-react";
import type React from "react";
import { useEffect, useState } from "react";
import type { config } from "../../wailsjs/go/models";
import "./Sidebar.css";

interface SidebarProps {
  items: config.Item[];
  onSelect: (item: config.Item) => void;
  onAdd: () => void;
  currentView: string;
  onViewChange: (view: any) => void;
  onCompare: (item1: config.Item, item2: config.Item) => void;
  selectedItem?: config.Item | null; // Add selectedItem prop to highlight correctly
}

const Sidebar: React.FC<SidebarProps> = ({
  items,
  onSelect,
  onAdd,
  currentView,
  onViewChange,
  selectedItem,
}) => {
  const [expandedProviders, setExpandedProviders] = useState<
    Record<string, boolean>
  >({});
  const [searchTerm, setSearchTerm] = useState("");

  // Group items by provider
  const groupedItems = items.reduce(
    (acc, item) => {
      const provider = item.provider || "Other";
      if (!acc[provider]) {
        acc[provider] = [];
      }
      acc[provider].push(item);
      return acc;
    },
    {} as Record<string, config.Item[]>,
  );

  // Initialize expanded state for all providers on first load if not set
  useEffect(() => {
    if (Object.keys(expandedProviders).length === 0 && items.length > 0) {
      const initialExpanded: Record<string, boolean> = {};
      Object.keys(groupedItems).forEach((p) => {
        initialExpanded[p] = true;
      });
      setExpandedProviders(initialExpanded);
    }
  }, [items, groupedItems, expandedProviders]);

  const toggleProvider = (provider: string) => {
    setExpandedProviders((prev) => ({
      ...prev,
      [provider]: !prev[provider],
    }));
  };

  const filteredProviders = Object.keys(groupedItems).filter((provider) => {
    const providerMatch = provider
      .toLowerCase()
      .includes(searchTerm.toLowerCase());
    const itemMatch = groupedItems[provider].some((item) =>
      item.name.toLowerCase().includes(searchTerm.toLowerCase()),
    );
    return providerMatch || itemMatch;
  });

  return (
    <aside className="sidebar">
      <div className="sidebar-header-section">
        <div className="app-brand">
          <div className="app-icon">
            <Terminal size={20} />
          </div>
          <h1 className="app-title">EasyConfig</h1>
        </div>
        <div className="sidebar-search">
          <Search size={16} />
          <input
            type="text"
            placeholder="Search..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
          />
        </div>
      </div>

      <nav className="sidebar-nav">
        <button
          type="button"
          className={`nav-item ${currentView === "health" ? "active" : ""}`}
          onClick={() => onViewChange("health")}
        >
          <Activity size={18} />
          <span>Health</span>
        </button>
        <button
          type="button"
          className={`nav-item ${currentView === "workflows" ? "active" : ""}`}
          onClick={() => onViewChange("workflows")}
        >
          <Workflow size={18} />
          <span>Workflows</span>
        </button>
        <button
          type="button"
          className={`nav-item ${currentView === "marketplace" ? "active" : ""}`}
          onClick={() => onViewChange("marketplace")}
        >
          <ShoppingBag size={18} />
          <span>Marketplace</span>
        </button>
        <button
          type="button"
          className={`nav-item ${currentView === "docs" ? "active" : ""}`}
          onClick={() => onViewChange("docs")}
        >
          <BookOpen size={18} />
          <span>Docs</span>
        </button>
      </nav>

      <div className="sidebar-divider" />

      <h3 className="sidebar-section-title">Configurations</h3>

      <div className="configs-list">
        {filteredProviders.map((provider) => (
          <div key={provider} className="provider-group">
            <button
              type="button"
              className="provider-header"
              onClick={() => toggleProvider(provider)}
            >
              {expandedProviders[provider] ? (
                <ChevronDown size={14} />
              ) : (
                <ChevronRight size={14} />
              )}
              <Brain size={14} />
              <span>{provider}</span>
            </button>
            {expandedProviders[provider] && (
              <div className="provider-items">
                {groupedItems[provider]
                  .filter((item) =>
                    item.name.toLowerCase().includes(searchTerm.toLowerCase()),
                  )
                  .map((item) => (
                    <button
                      key={item.path}
                      type="button"
                      className={`config-item ${
                        selectedItem?.path === item.path &&
                        currentView === "configs"
                          ? "active"
                          : ""
                      }`}
                      onClick={() => {
                        onSelect(item);
                        onViewChange("configs");
                      }}
                    >
                      <FileJson size={14} />
                      <span>{item.name}</span>
                    </button>
                  ))}
              </div>
            )}
          </div>
        ))}
      </div>

      <div className="sidebar-footer">
        <button type="button" className="btn-add-config" onClick={onAdd}>
          <Plus size={16} />
          <span>Add Configuration</span>
        </button>
      </div>
    </aside>
  );
};

export default Sidebar;
