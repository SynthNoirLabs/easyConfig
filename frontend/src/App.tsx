import { useState } from "react";
import { Toaster, toast } from "sonner";
import "./App.css";
import type { config } from "../wailsjs/go/models";
import AddConfigModal from "./components/AddConfigModal";
import ConfigEditor from "./components/ConfigEditor";
import Docs from "./components/Docs";
import Layout from "./components/Layout";
import Marketplace from "./components/Marketplace";
import Sidebar from "./components/Sidebar";
import Workflows from "./components/Workflows";
import { useConfig } from "./context/ConfigContext";

function AppContent() {
  const { configs, loading, error, refreshConfigs } = useConfig();
  const [selectedItem, setSelectedItem] = useState<config.Item | null>(null);
  const [isAddModalOpen, setIsAddModalOpen] = useState(false);
  const [currentView, setCurrentView] = useState<
    "configs" | "workflows" | "marketplace" | "docs"
  >("configs");

  const handleSelectConfig = (item: config.Item) => {
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

  const renderContent = () => {
    switch (currentView) {
      case "workflows":
        return <Workflows />;
      case "marketplace":
        return <Marketplace />; // We need to import this
      case "docs":
        return <Docs />;
      default:
        return selectedItem ? (
          <ConfigEditor configItem={selectedItem} />
        ) : (
          <div className="empty-state">
            <h2>Welcome to easyConfig</h2>
            <p>
              Select a configuration file to edit, or explore workflows and
              marketplace.
            </p>
          </div>
        );
    }
  };

  return (
    <>
      <Layout
        sidebar={
          <Sidebar
            items={configs}
            onSelect={handleSelectConfig}
            onAdd={handleOpenAddModal}
            currentView={currentView}
            onViewChange={setCurrentView}
          />
        }
      >
        <div className="app-content">{renderContent()}</div>
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
