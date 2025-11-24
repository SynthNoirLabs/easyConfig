import { useState } from 'react';
import './App.css';
import Layout from './components/Layout';
import Sidebar, { ConfigItem } from './components/Sidebar';

// Dummy data for testing
const dummyConfigItems: ConfigItem[] = [
  {
    name: 'claude_desktop_config.json',
    provider: 'Claude Code',
    path: '~/.config/claude/config.json',
  },
  {
    name: 'mcp_servers.json',
    provider: 'Claude Code',
    path: '~/.config/claude/mcp_servers.json',
  },
  {
    name: 'hooks.json',
    provider: 'Claude Code',
    path: '~/.config/claude/hooks.json',
  },
  {
    name: 'settings.json',
    provider: 'VS Code',
    path: '~/.config/Code/User/settings.json',
  },
  {
    name: 'keybindings.json',
    provider: 'VS Code',
    path: '~/.config/Code/User/keybindings.json',
  },
  {
    name: 'gemini_config.json',
    provider: 'Gemini',
    path: '~/.config/gemini/config.json',
  },
  {
    name: 'model_settings.json',
    provider: 'Gemini',
    path: '~/.config/gemini/model_settings.json',
  },
];

function App() {
  const [selectedItem, setSelectedItem] = useState<ConfigItem | null>(null);

  const handleSelectConfig = (item: ConfigItem) => {
    setSelectedItem(item);
  };

  return (
    <Layout
      sidebar={
        <Sidebar
          items={dummyConfigItems}
          onSelect={handleSelectConfig}
        />
      }
    >
      <div className="app-content">
        {selectedItem ? (
          <div className="config-details">
            <h1 className="config-title">{selectedItem.name}</h1>
            <div className="config-meta">
              <div className="config-meta-item">
                <span className="config-meta-label">Provider:</span>
                <span className="config-meta-value">{selectedItem.provider}</span>
              </div>
              <div className="config-meta-item">
                <span className="config-meta-label">Path:</span>
                <span className="config-meta-value">{selectedItem.path}</span>
              </div>
            </div>
            <div className="config-placeholder">
              <p>Configuration editor will be displayed here.</p>
              <p className="config-placeholder-note">
                This is where you'll be able to view and edit the configuration file content.
              </p>
            </div>
          </div>
        ) : (
          <div className="empty-state">
            <h2>Welcome to easyConfig</h2>
            <p>Select a configuration file from the sidebar to get started.</p>
          </div>
        )}
      </div>
    </Layout>
  );
}

export default App;
