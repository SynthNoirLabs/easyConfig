import { Toaster, toast } from "sonner";
import { useState } from "react";
import "./App.css";
import type { config } from "../wailsjs/go/config/models";
import AddConfigModal from "./components/AddConfigModal";
import ConfigEditor from "./components/ConfigEditor";
import Layout from "./components/Layout";
import Sidebar from "./components/Sidebar";
import { useConfig } from "./context/ConfigContext";

function AppContent() {
  const { configs, loading, error, refreshConfigs } = useConfig();
  const [selectedItem, setSelectedItem] = useState<config.ConfigItem | null>(
    null,
  );
  const [isAddModalOpen, setIsAddModalOpen] = useState(false);

  const handleSelectConfig = (item: config.ConfigItem) => {
    setSelectedItem(item);
  };

  const handleOpenAddModal = () => {
    setIsAddModalOpen(true);
  };

  const handleCloseAddModal = () => {
    setIsAddModalOpen(false);
  };

  const handleConfigAdded = async () => {
    await refreshConfigs();
    toast.success("Configuration created successfully");
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
    <>
      <Layout
        sidebar={
          <Sidebar
            items={configs}
            onSelect={handleSelectConfig}
            onAdd={handleOpenAddModal}
          />
        }
      >
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
      <AddConfigModal
        isOpen={isAddModalOpen}
        onClose={handleCloseAddModal}
        onSuccess={handleConfigAdded}
      />
    </>
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
