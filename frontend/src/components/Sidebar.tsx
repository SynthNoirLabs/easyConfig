import { Box, FileJson, Plus, Trash2 } from "lucide-react";
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
}

const Sidebar: React.FC<SidebarProps> = ({ items, onSelect, onAdd }) => {
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
          // Optionally clear editor content via parent callback, but simpler to just deselect
        }
      } catch (err) {
        toast.error("Failed to delete configuration");
      }
    }
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
                <div
                  key={item.path}
                  className={`sidebar-item ${selectedPath === item.path ? "sidebar-item-active" : ""}`}
                  onClick={() => handleItemClick(item)}
                  onKeyDown={(e) => {
                    if (e.key === "Enter" || e.key === " ") {
                      handleItemClick(item);
                    }
                  }}
                  role="button"
                  tabIndex={0}
                >
                  <FileJson size={16} className="sidebar-item-icon" />
                  <div className="sidebar-item-content">
                    <div className="sidebar-item-name">{item.name}</div>
                    <div className="sidebar-item-path">{item.path}</div>
                  </div>
                  <button
                    type="button"
                    className="btn-delete"
                    onClick={(e) => handleDelete(e, item)}
                    title="Delete"
                  >
                    <Trash2 size={14} />
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
