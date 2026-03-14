import { useEffect, useState } from "react";
import { Toaster, toast } from "sonner";
import "./App.css";
import { Settings } from "lucide-react";
import type { config } from "../wailsjs/go/models";
import AddConfigModal from "./components/AddConfigModal";
import CommandPalette from "./components/CommandPalette";
import ComparisonViewer from "./components/ComparisonViewer";
import ConfigEditor from "./components/ConfigEditor";
import ConfigWizard from "./components/ConfigWizard";
import Docs from "./components/Docs";
import ErrorBoundary from "./components/ErrorBoundary";
import HealthDashboard from "./components/HealthDashboard";
import Layout from "./components/Layout";
import Marketplace from "./components/Marketplace";
import ShortcutsModal from "./components/ShortcutsModal";
import Sidebar from "./components/Sidebar";
import { SettingsPage } from "./components/SettingsPage";
import Workflows from "./components/Workflows";
import { useConfig } from "./context/ConfigContext";
import { useKeyboardShortcuts } from "./hooks/useKeyboardShortcuts";

type SelectableItem = config.Item & { initialLine?: number };

function AppContent() {
  const { configs, loading, error, refreshConfigs } = useConfig();
  const [selectedItem, setSelectedItem] = useState<SelectableItem | null>(null);
  const [comparisonItems, setComparisonItems] = useState<
    [config.Item, config.Item] | null
  >(null);
  const [isAddModalOpen, setIsAddModalOpen] = useState(false);
  const [isShortcutsModalOpen, setIsShortcutsModalOpen] = useState(false);
  const [isCommandPaletteOpen, setIsCommandPaletteOpen] = useState(false);
  const [currentView, setCurrentView] = useState<
    "configs" | "health" | "workflows" | "marketplace" | "docs" | "settings"
  >("configs");

  const handleSelectConfig = (item: config.Item) => {
    setSelectedItem(item);
    setComparisonItems(null); // Exit comparison mode when a single item is selected
    setCurrentView("configs");
  };

  const handleCompareConfigs = (item1: config.Item, item2: config.Item) => {
    setComparisonItems([item1, item2]);
    setSelectedItem(null); // Deselect single item view
  };

  const handleSearch = (path: string, line?: number) => {
    const item = configs.find((c) => c.path === path);
    if (item) {
      setSelectedItem({ ...item, initialLine: line });
      setCurrentView("configs");
    }
    setIsCommandPaletteOpen(false);
  };

  const handleOpenAddModal = () => {
    setIsAddModalOpen(true);
  };

  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if ((e.ctrlKey || e.metaKey) && e.key === "k") {
        e.preventDefault();
        setIsCommandPaletteOpen(true);
      }
    };
    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, []);

  const handleCloseAddModal = () => {
    setIsAddModalOpen(false);
  };

  const handleConfigAdded = async () => {
    await refreshConfigs();
    toast.success("Configuration created successfully");
  };

  useKeyboardShortcuts({
    "ctrl+n": handleOpenAddModal,
    "?": () => setIsShortcutsModalOpen((prev) => !prev),
  });

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

  const getHeaderTitle = () => {
    if (comparisonItems) return "Compare Configurations";
    switch (currentView) {
      case "health":
        return "System Health";
      case "workflows":
        return "Workflows";
      case "docs":
        return "Documentation";
      case "marketplace":
        return "Marketplace";
      default:
        return selectedItem ? selectedItem.name : "Mission Control";
    }
  };

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
      case "health":
        return <HealthDashboard />;
      case "workflows":
        return <Workflows />;
      case "docs":
        return <Docs />;
      case "marketplace":
        return <Marketplace />;
      case "settings":
        return <SettingsPage />;
      default:
        return selectedItem ? (
          <ConfigEditor configItem={selectedItem} />
        ) : (
          <div className="empty-state">
            <h2>Welcome to EasyConfig</h2>
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
              selectedItem={selectedItem}
            />
          </ErrorBoundary>
        }
      >
        <div className="app-content">
          <header className="app-header">
            <div className="header-left">
              <h2 className="header-title">{getHeaderTitle()}</h2>
              {currentView === "configs" && selectedItem && (
                <span className="badge badge-active">Active</span>
              )}
            </div>
            <div className="header-actions">
              <button type="button" className="btn-icon" title="Settings">
                <Settings size={18} />
              </button>
            </div>
          </header>
          <div className="main-scroll-area">
            <ErrorBoundary>{renderContent()}</ErrorBoundary>
          </div>
        </div>
      </Layout>

      <AddConfigModal
        isOpen={isAddModalOpen}
        onClose={handleCloseAddModal}
        onSuccess={handleConfigAdded}
      />

      <ShortcutsModal
        isOpen={isShortcutsModalOpen}
        onClose={() => setIsShortcutsModalOpen(false)}
      />

      {isCommandPaletteOpen && (
        <CommandPalette
          onSelect={handleSearch}
          onClose={() => setIsCommandPaletteOpen(false)}
        />
      )}
    </>
  );
}

function App() {
  return (
    <>
      <Toaster richColors theme="dark" />
      <AppContent />
    </>
  );
}

export default App;
