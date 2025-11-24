import { useState } from 'react';
import './App.css';
import Layout from './components/Layout';
import Sidebar from './components/Sidebar';
import { useConfig } from './context/ConfigContext';
import { config } from '../wailsjs/go/config/models';

function App() {
  const { configs, loading, error } = useConfig();
  const [selectedItem, setSelectedItem] = useState<config.ConfigItem | null>(null);

  const handleSelectConfig = (item: config.ConfigItem) => {
    setSelectedItem(item);
  };

  if (loading) {
    return <div className="loading">Loading configurations...</div>;
  }

  if (error) {
    return <div className="error">Error: {error}</div>;
  }

  return (
    <Layout
      sidebar={
        <Sidebar
          items={configs}
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
