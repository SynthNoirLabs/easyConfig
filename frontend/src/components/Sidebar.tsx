import { Box, FileJson, Plus } from "lucide-react";
import type React from "react";
import { useState } from "react";
import type { config } from "../../wailsjs/go/config/models";
import "./Sidebar.css";

interface SidebarProps {
  items: config.ConfigItem[];
  onSelect: (item: config.ConfigItem) => void;
  onAdd: () => void;
}

const Sidebar: React.FC<SidebarProps> = ({ items, onSelect, onAdd }) => {
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
    onSelect(item);
  };

  return (
    <div className="sidebar">
      <div className="sidebar-header">
        <h2 className="sidebar-title">Configuration Files</h2>
        <button className="btn-add" onClick={onAdd} title="Add Configuration">
          <Plus size={18} />
        </button>
      </div>
      <div className="sidebar-content">
        {Object.entries(groupedItems).map(([provider, providerItems]) => (
          <div key={provider} className="sidebar-group">
            <div className="sidebar-group-header">
              <Box size={16} className="sidebar-group-icon" />
              <span className="sidebar-group-title">{provider}</span>
            </div>
            <div className="sidebar-group-items">
              {providerItems.map((item) => (
                <button
                  key={item.path}
                  className={`sidebar-item ${selectedPath === item.path ? "sidebar-item-active" : ""}`}
                  onClick={() => handleItemClick(item)}
                  onKeyDown={(e) => {
                    if (e.key === "Enter" || e.key === " ") {
                      handleItemClick(item);
                    }
                  }}
                  type="button"
                >
                  <FileJson size={16} className="sidebar-item-icon" />
                  <div className="sidebar-item-content">
                    <div className="sidebar-item-name">{item.name}</div>
                    <div className="sidebar-item-path">{item.path}</div>
                  </div>
                </button>
              ))}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default Sidebar;
