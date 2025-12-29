import { useState } from "react";
import { Toaster, toast } from "sonner";
import "./App.css";
import type { config } from "../wailsjs/go/models";
import AddConfigModal from "./components/AddConfigModal";
import ComparisonViewer from "./components/ComparisonViewer";
import ConfigEditor from "./components/ConfigEditor";
import ConfigWizard from "./components/ConfigWizard";
import Docs from "./components/Docs";
import Layout from "./components/Layout";
import Marketplace from "./components/Marketplace";
import ErrorBoundary from "./components/ErrorBoundary";
import Sidebar from "./components/Sidebar";
import Workflows from "./components/Workflows";
import { useConfig } from "./context/ConfigContext";

function AppContent() {
  const { configs, loading, error, refreshConfigs } = useConfig();
  const [selectedItem, setSelectedItem] = useState<config.Item | null>(null);
  const [comparisonItems, setComparisonItems] = useState<
    [config.Item, config.Item] | null
  >(null);
  const [isAddModalOpen, setIsAddModalOpen] = useState(false);
  const [currentView, setCurrentView] = useState<
    "configs" | "workflows" | "marketplace" | "docs"
  >("configs");

  const handleSelectConfig = (item: config.Item) => {
    setSelectedItem(item);
    setComparisonItems(null); // Exit comparison mode when a single item is selected
  };

  const handleCompareConfigs = (item1: config.Item, item2: config.Item) => {
    setComparisonItems([item1, item2]);
    setSelectedItem(null); // Deselect single item view
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
    if (comparisonItems) {
      return (
        <ComparisonViewer
          item1={comparisonItems[0]}
          item2={comparisonItems[1]}
          onClose={() => setComparisonItems(null)}
        />
      );
    }

    switch (currentView) {
      case "workflows":
        return <Workflows />;
      case "docs":
        return <Docs />;
      case "marketplace":
        return <Marketplace />;
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
            <ConfigWizard />
          </div>
        );
    }
  };

  return (
    <>
      <Layout
        sidebar={
          <ErrorBoundary>
            <Sidebar
              items={configs}
              onSelect={handleSelectConfig}
              onAdd={handleOpenAddModal}
              currentView={currentView}
              onViewChange={setCurrentView}
              onCompare={handleCompareConfigs}
            />
          </ErrorBoundary>
        }
      >
        <div className="app-content">
          <ErrorBoundary>{renderContent()}</ErrorBoundary>
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
