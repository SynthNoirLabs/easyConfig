import { Toaster, toast } from "sonner";
import { useState } from "react";
import "./App.css";
import type { config } from "../wailsjs/go/config/models";
import AddConfigModal from "./components/AddConfigModal";
import ConfigEditor from "./components/ConfigEditor";
import Layout from "./components/Layout";
import Sidebar from "./components/Sidebar";
import Workflows from "./components/Workflows";
import { useConfig } from "./context/ConfigContext";

function AppContent() {
  const { configs, loading, error, refreshConfigs } = useConfig();
  const [selectedItem, setSelectedItem] = useState<config.ConfigItem | null>(null);
  const [isAddModalOpen, setIsAddModalOpen] = useState(false);
  const [currentView, setCurrentView] = useState<"configs" | "workflows">("configs");

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
      <div className="app-nav">
        <button 
          className={`nav-item ${currentView === "configs" ? "active" : ""}`}
          onClick={() => setCurrentView("configs")}
        >
          Configs
        </button>
        <button 
          className={`nav-item ${currentView === "workflows" ? "active" : ""}`}
          onClick={() => setCurrentView("workflows")}
        >
          Workflows
        </button>
      </div>

      {currentView === "configs" ? (
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
      ) : (
        <Layout sidebar={null}>
           <Workflows />
        </Layout>
      )}

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
