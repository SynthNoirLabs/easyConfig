import { Toaster, toast } from "sonner";
import { useState } from "react";
import "./App.css";
import { config } from "../wailsjs/go/models";
import AddConfigModal from "./components/AddConfigModal";
import ConfigEditor from "./components/ConfigEditor";
import Layout from "./components/Layout";
import Sidebar from "./components/Sidebar";
import Marketplace from "./components/Marketplace";
import { useConfig } from "./context/ConfigContext";
import { ShoppingBag } from "lucide-react";

function AppContent() {
  const { configs, loading, error, refreshConfigs } = useConfig();
  const [selectedItem, setSelectedItem] = useState<config.Item | null>(
    null,
  );
  const [isAddModalOpen, setIsAddModalOpen] = useState(false);
  const [view, setView] = useState<"editor" | "marketplace">("editor");

  const handleSelectConfig = (item: config.Item) => {
    setSelectedItem(item);
    setView("editor");
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
          <div className="sidebar-container">
            <div className="sidebar-nav">
               <button 
                 className={`nav-btn ${view === "marketplace" ? "active" : ""}`}
                 onClick={() => setView("marketplace")}
               >
                 <ShoppingBag size={16} /> Marketplace
               </button>
            </div>
            <Sidebar
              items={configs}
              onSelect={handleSelectConfig}
              onAdd={handleOpenAddModal}
            />
          </div>
        }
      >
        <div className="app-content">
          {view === "marketplace" ? (
            <Marketplace />
          ) : selectedItem ? (
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
