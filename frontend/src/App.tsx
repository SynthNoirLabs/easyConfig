import { Toaster } from "sonner";
import { useState } from "react";
import "./App.css";
import type { config } from "../wailsjs/go/config/models";
import ConfigEditor from "./components/ConfigEditor";
import Layout from "./components/Layout";
import Sidebar from "./components/Sidebar";
import { useConfig } from "./context/ConfigContext";

function AppContent() {
  const { configs, loading, error } = useConfig();
  const [selectedItem, setSelectedItem] = useState<config.ConfigItem | null>(
    null,
  );

  const handleSelectConfig = (item: config.ConfigItem) => {
    setSelectedItem(item);
  };

  if (loading) {
    return (
      <div className="app-loading">
        <p>Loading configurations...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="app-error">
        <p>Error: {error}</p>
      </div>
    );
  }

  return (
    <Layout sidebar={<Sidebar items={configs} onSelect={handleSelectConfig} />}>
      <div className="app-content">
        {selectedItem ? (
          <ConfigEditor configItem={selectedItem} />
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

function App() {
  return (
    <>
      <Toaster richColors />
      <AppContent />
    </>
  );
}

export default App;
