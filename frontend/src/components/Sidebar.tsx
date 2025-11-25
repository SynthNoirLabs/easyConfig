import { Box, FileJson, Plus, Trash2, LayoutGrid, Workflow, Store } from "lucide-react";
import type React from "react";
import { useState } from "react";
import { toast } from "sonner";
import type { config } from "../../wailsjs/go/config/models";
import { useConfig } from "../context/ConfigContext";
import "./Sidebar.css";

interface SidebarProps {
  items: config.ConfigItem[];
  onSelect: (item: config.ConfigItem) => void;
  onAdd: () => void;
  currentView: "configs" | "workflows" | "marketplace";
  onViewChange: (view: "configs" | "workflows" | "marketplace") => void;
}

const Sidebar: React.FC<SidebarProps> = ({ items, onSelect, onAdd, currentView, onViewChange }) => {
  const { deleteConfig } = useConfig();
  const [selectedPath, setSelectedPath] = useState<string | null>(null);

  // Group items by provider
  const groupedItems = items.reduce(
    (acc, item) => {
      if (!acc[item.provider]) {
        acc[item.provider] = [];
      }
      acc[item.provider].push(item);
      return acc;
    },
    {} as Record<string, config.ConfigItem[]>,
  );

  const handleItemClick = (item: config.ConfigItem) => {
    setSelectedPath(item.path);
    onViewChange("configs"); // Switch to configs view when selecting a file
    onSelect(item);
  };

  const handleDelete = async (e: React.MouseEvent, item: config.ConfigItem) => {
    e.stopPropagation();
    if (confirm(`Are you sure you want to delete ${item.name}?`)) {
      try {
        await deleteConfig(item.path);
        toast.success("Configuration deleted");
        if (selectedPath === item.path) {
          setSelectedPath(null);
        }
      } catch (err) {
        toast.error("Failed to delete configuration");
      }
    }
  };

  return (
    <div className="sidebar">
      <div className="sidebar-section">
        <h3 className="sidebar-section-title">Menu</h3>
        <div 
          className={`sidebar-nav-item ${currentView === "configs" && !selectedPath ? "active" : ""}`}
          onClick={() => {
            onViewChange("configs");
            setSelectedPath(null); // Deselect specific file to show dashboard/empty state
          }}
        >
          <LayoutGrid size={18} />
          <span>Dashboard</span>
        </div>
        <div 
          className={`sidebar-nav-item ${currentView === "workflows" ? "active" : ""}`}
          onClick={() => onViewChange("workflows")}
        >
          <Workflow size={18} />
          <span>Workflows</span>
        </div>
        <div 
          className={`sidebar-nav-item ${currentView === "marketplace" ? "active" : ""}`}
          onClick={() => onViewChange("marketplace")}
        >
          <Store size={18} />
          <span>Marketplace</span>
        </div>
      </div>

      <div className="sidebar-divider" />

      <div className="sidebar-header">
        <h2 className="sidebar-title">Configurations</h2>
        <button className="btn-icon" onClick={onAdd} title="Add Configuration">
          <Plus size={16} />
        </button>
      </div>
      
      <div className="sidebar-content">
        {Object.entries(groupedItems).map(([provider, providerItems]) => (
          <div key={provider} className="sidebar-group">
            <div className="sidebar-group-header">
              <span className="sidebar-group-title">{provider}</span>
            </div>
            <div className="sidebar-group-items">
              {providerItems.map((item) => (
                <div
                  key={item.path}
                  className={`sidebar-item ${selectedPath === item.path && currentView === "configs" ? "sidebar-item-active" : ""}`}
                  onClick={() => handleItemClick(item)}
                  role="button"
                  tabIndex={0}
                >
                  <FileJson size={14} className="sidebar-item-icon" />
                  <span className="sidebar-item-name">{item.name}</span>
                  <button
                    type="button"
                    className="btn-delete-icon"
                    onClick={(e) => handleDelete(e, item)}
                    title="Delete"
                  >
                    <Trash2 size={12} />
                  </button>
                </div>
              ))}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default Sidebar;
