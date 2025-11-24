import React, { useState } from 'react';
import { FileJson, Box } from 'lucide-react';
import './Sidebar.css';

export interface ConfigItem {
  name: string;
  provider: string;
  path: string;
}

interface SidebarProps {
  items: ConfigItem[];
  onSelect: (item: ConfigItem) => void;
}

const Sidebar: React.FC<SidebarProps> = ({ items, onSelect }) => {
  const [selectedPath, setSelectedPath] = useState<string | null>(null);

  // Group items by provider
  const groupedItems = items.reduce((acc, item) => {
    if (!acc[item.provider]) {
      acc[item.provider] = [];
    }
    acc[item.provider].push(item);
    return acc;
  }, {} as Record<string, ConfigItem[]>);

  const handleItemClick = (item: ConfigItem) => {
    setSelectedPath(item.path);
    onSelect(item);
  };

  return (
    <div className="sidebar">
      <div className="sidebar-header">
        <h2 className="sidebar-title">Configuration Files</h2>
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
                  className={`sidebar-item ${selectedPath === item.path ? 'sidebar-item-active' : ''}`}
                  onClick={() => handleItemClick(item)}
                >
                  <FileJson size={16} className="sidebar-item-icon" />
                  <div className="sidebar-item-content">
                    <div className="sidebar-item-name">{item.name}</div>
                    <div className="sidebar-item-path">{item.path}</div>
                  </div>
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
