import {
  BookOpen,
  Clipboard,
  Copy,
  ExternalLink,
  FileJson,
  LayoutGrid,
  Plus,
  Store,
  Trash2,
  Workflow,
} from "lucide-react";
import type React from "react";
import { useCallback, useEffect, useState } from "react";
import { toast } from "sonner";
import {
  ApplyProfile,
  DeleteProfile,
  ListProfiles,
  SaveProfile,
} from "../../wailsjs/go/main/App";
import type { config, config as configModels } from "../../wailsjs/go/models";
import { BrowserOpenURL } from "../../wailsjs/runtime/runtime";
import { useConfig } from "../context/ConfigContext";
import "./Sidebar.css";

interface SidebarProps {
  items: config.Item[];
  onSelect: (item: config.Item) => void;
  onAdd: () => void;
  currentView: "configs" | "workflows" | "marketplace" | "docs";
  onViewChange: (
    view: "configs" | "workflows" | "marketplace" | "docs",
  ) => void;
}

const Sidebar: React.FC<SidebarProps> = ({
  items,
  onSelect,
  onAdd,
  currentView,
  onViewChange,
}) => {
  const { deleteConfig } = useConfig();
  const { readConfig } = useConfig();
  const [selectedPath, setSelectedPath] = useState<string | null>(null);
  const [profiles, setProfiles] = useState<configModels.ProfileSummary[]>([]);
  const [selectedProfile, setSelectedProfile] = useState<string>("");

  const refreshProfiles = useCallback(async () => {
    try {
      const data = await ListProfiles();
      setProfiles(data);
      if (data.length && !selectedProfile) {
        setSelectedProfile(data[0].name);
      }
    } catch (_err) {
      // silent
    }
  }, [selectedProfile]);

  // Lazy load once
  useEffect(() => {
    void refreshProfiles();
  }, [refreshProfiles]);

  const scopeOrder: Array<"global" | "project" | "system"> = [
    "global",
    "project",
    "system",
  ];
  const scopeLabels: Record<string, string> = {
    global: "Global",
    project: "Project",
    system: "System",
  };

  const groupedByScope = scopeOrder
    .map((scope) => {
      const inScope = items.filter((item) => item.scope === scope);
      const byProvider = inScope.reduce(
        (acc, item) => {
          if (!acc[item.provider]) acc[item.provider] = [];
          acc[item.provider].push(item);
          return acc;
        },
        {} as Record<string, config.Item[]>,
      );
      return { scope, providers: byProvider };
    })
    .filter((group) => Object.values(group.providers).flat().length > 0);

  const handleItemClick = (item: config.Item) => {
    setSelectedPath(item.path);
    onViewChange("configs"); // Switch to configs view when selecting a file
    onSelect(item);
  };

  const handleDelete = async (e: React.MouseEvent, item: config.Item) => {
    e.stopPropagation();
    if (confirm(`Are you sure you want to delete ${item.name}?`)) {
      try {
        await deleteConfig(item.path);
        toast.success("Configuration deleted");
        if (selectedPath === item.path) {
          setSelectedPath(null);
        }
      } catch (_err) {
        toast.error("Failed to delete configuration");
      }
    }
  };

  const handleCopyPath = async (
    e: React.MouseEvent,
    item: config.Item,
  ): Promise<void> => {
    e.stopPropagation();
    try {
      await navigator.clipboard.writeText(item.path);
      toast.success("Path copied");
    } catch (err) {
      console.error(err);
      toast.error("Could not copy path");
    }
  };

  const handleCopyContent = async (
    e: React.MouseEvent,
    item: config.Item,
  ): Promise<void> => {
    e.stopPropagation();
    try {
      const content = await readConfig(item.path);
      await navigator.clipboard.writeText(content);
      toast.success("Content copied");
    } catch (err) {
      console.error(err);
      toast.error("Could not copy content");
    }
  };

  const handleOpenExternal = async (
    e: React.MouseEvent,
    item: config.Item,
  ): Promise<void> => {
    e.stopPropagation();
    try {
      await BrowserOpenURL(`file://${item.path}`);
    } catch (err) {
      console.error(err);
      toast.error("Could not open file");
    }
  };

  const handleApplyProfile = async () => {
    if (!selectedProfile) return;
    if (
      !confirm(
        `Apply profile ${selectedProfile}? This overwrites existing configs.`,
      )
    ) {
      return;
    }
    try {
      const written = await ApplyProfile(selectedProfile);
      toast.success(
        `Applied profile ${selectedProfile}${written.length ? ` (${written.length} files)` : ""}`,
      );
      setSelectedPath(null);
    } catch (err) {
      toast.error(`Failed to apply profile: ${err}`);
    }
  };

  const handleSaveProfile = async () => {
    const name = prompt(
      "Profile name (letters, numbers, - _ .):",
      "my-profile",
    );
    if (!name) return;
    try {
      await SaveProfile(name);
      toast.success(`Saved profile ${name}`);
      await refreshProfiles();
      setSelectedProfile(name);
    } catch (err) {
      toast.error(`Failed to save profile: ${err}`);
    }
  };

  const handleDeleteProfile = async () => {
    if (!selectedProfile) return;
    if (!confirm(`Delete profile ${selectedProfile}?`)) return;
    try {
      await DeleteProfile(selectedProfile);
      toast.success(`Deleted profile ${selectedProfile}`);
      setSelectedProfile("");
      await refreshProfiles();
    } catch (err) {
      toast.error(`Failed to delete profile: ${err}`);
    }
  };

  return (
    <div className="sidebar">
      <div className="sidebar-section">
        <h3 className="sidebar-section-title">Menu</h3>
        <button
          type="button"
          className={`sidebar-nav-item ${currentView === "configs" && !selectedPath ? "active" : ""}`}
          onClick={() => {
            onViewChange("configs");
            setSelectedPath(null); // Deselect specific file to show dashboard/empty state
          }}
        >
          <LayoutGrid size={18} />
          <span>Dashboard</span>
        </button>
        <button
          type="button"
          className={`sidebar-nav-item ${currentView === "workflows" ? "active" : ""}`}
          onClick={() => onViewChange("workflows")}
        >
          <Workflow size={18} />
          <span>Workflows</span>
        </button>
        <button
          type="button"
          className={`sidebar-nav-item ${currentView === "marketplace" ? "active" : ""}`}
          onClick={() => onViewChange("marketplace")}
        >
          <Store size={18} />
          <span>Marketplace</span>
        </button>
        <button
          type="button"
          className={`sidebar-nav-item ${currentView === "docs" ? "active" : ""}`}
          onClick={() => onViewChange("docs")}
        >
          <BookOpen size={18} />
          <span>Docs</span>
        </button>
      </div>

      <div className="sidebar-divider" />

      <div className="sidebar-header">
        <h2 className="sidebar-title">Configurations</h2>
        <button
          type="button"
          className="btn-icon"
          onClick={onAdd}
          title="Add Configuration"
        >
          <Plus size={16} />
        </button>
      </div>

      <div className="profiles-bar">
        <select
          className="profile-select"
          value={selectedProfile}
          onChange={(e) => setSelectedProfile(e.target.value)}
        >
          <option value="">Select profile</option>
          {profiles.map((p) => (
            <option key={p.name} value={p.name}>
              {p.name} ({p.itemCount})
            </option>
          ))}
        </select>
        <div className="profile-actions">
          <button
            type="button"
            className="btn-secondary"
            onClick={handleApplyProfile}
            disabled={!selectedProfile}
          >
            Apply
          </button>
          <button
            type="button"
            className="btn-secondary"
            onClick={handleSaveProfile}
          >
            Save current
          </button>
          <button
            type="button"
            className="btn-secondary"
            onClick={handleDeleteProfile}
            disabled={!selectedProfile}
          >
            Delete
          </button>
        </div>
      </div>

      <div className="sidebar-content">
        {groupedByScope.map(({ scope, providers }) => (
          <div key={scope} className="sidebar-scope">
            <div className="sidebar-scope-header">
              <span className="sidebar-scope-title">
                {scopeLabels[scope] ?? scope}
              </span>
            </div>
            {Object.entries(providers).map(([provider, providerItems]) => (
              <div key={`${scope}-${provider}`} className="sidebar-group">
                <div className="sidebar-group-header">
                  <span className="sidebar-group-title">{provider}</span>
                </div>
                <div className="sidebar-group-items">
                  {providerItems.map((item) => (
                    <div
                      key={item.path}
                      className={`sidebar-item ${selectedPath === item.path && currentView === "configs" ? "sidebar-item-active" : ""}`}
                    >
                      <button
                        type="button"
                        className="sidebar-item-main"
                        onClick={() => handleItemClick(item)}
                      >
                        <FileJson size={14} className="sidebar-item-icon" />
                        <span className="sidebar-item-name">{item.name}</span>
                      </button>
                      <div className="sidebar-item-actions">
                        <button
                          type="button"
                          className="btn-ghost"
                          onClick={(e) => handleCopyPath(e, item)}
                          title="Copy path"
                        >
                          <Copy size={14} />
                        </button>
                        <button
                          type="button"
                          className="btn-ghost"
                          onClick={(e) => handleCopyContent(e, item)}
                          title="Copy content"
                        >
                          <Clipboard size={14} />
                        </button>
                        <button
                          type="button"
                          className="btn-ghost"
                          onClick={(e) => handleOpenExternal(e, item)}
                          title="Open externally"
                        >
                          <ExternalLink size={14} />
                        </button>
                      </div>
                      <button
                        type="button"
                        className="btn-delete-icon"
                        onClick={(e) => handleDelete(e, item)}
                        title="Delete"
                      >
                        <Trash2 size={12} />
                      </button>
                    </div>
                  ))}
                </div>
              </div>
            ))}
          </div>
        ))}
      </div>
    </div>
  );
};

export default Sidebar;
