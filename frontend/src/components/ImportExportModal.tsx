import type React from "react";
import { useEffect, useState } from "react";
import { toast } from "sonner";
import {
  ExportAllProfiles,
  ExportProfiles,
  ImportProfilesFromFile,
  ImportProfilesFromURL,
  ListProfiles,
} from "../../wailsjs/go/main/App";
import type { config } from "../../wailsjs/go/models";
import { OpenFile, SaveFile } from "../../wailsjs/runtime/runtime";
import "./ImportExportModal.css";

interface ImportExportModalProps {
  isOpen: boolean;
  onClose: () => void;
}

const ImportExportModal: React.FC<ImportExportModalProps> = ({
  isOpen,
  onClose,
}) => {
  const [activeTab, setActiveTab] = useState<"import" | "export">("import");
  const [profiles, setProfiles] = useState<config.ProfileSummary[]>([]);
  const [selectedProfiles, setSelectedProfiles] = useState<string[]>([]);
  const [importUrl, setImportUrl] = useState("");
  const [importStrategy, setImportStrategy] = useState("skip");

  useEffect(() => {
    if (isOpen) {
      ListProfiles()
        .then(setProfiles)
        .catch(() => toast.error("Failed to list profiles"));
    }
  }, [isOpen]);

  const handleExport = async (all = false) => {
    try {
      const names = all ? profiles.map((p) => p.name) : selectedProfiles;
      if (names.length === 0) {
        toast.warning("No profiles selected for export");
        return;
      }
      const data = await (all ? ExportAllProfiles() : ExportProfiles(names));
      const path = await SaveFile({
        defaultFilename: "profiles.easyconfig",
        filters: [{ displayName: "EasyConfig Profiles", pattern: "*.easyconfig" }],
      });
      if (path) {
        // Wails SaveFile doesn't write content, we just get a path
        // We need a way to write the data to the path. This is a limitation.
        // For now, we'll copy to clipboard and notify the user.
        await navigator.clipboard.writeText(new TextDecoder().decode(data));
        toast.success(
          "Export data copied to clipboard. Please save it to the selected file.",
        );
      }
    } catch (err) {
      toast.error(`Export failed: ${err}`);
    }
  };

  const handleImportFile = async () => {
    try {
      const path = await OpenFile({
        filters: [{ displayName: "EasyConfig Profiles", pattern: "*.easyconfig" }],
      });
      if (path) {
        const results = await ImportProfilesFromFile(path, importStrategy);
        toast.success(`Import complete: ${results.length} profiles processed.`);
        onClose();
      }
    } catch (err) {
      toast.error(`Import failed: ${err}`);
    }
  };

  const handleImportUrl = async () => {
    if (!importUrl) {
      toast.warning("Please enter a URL");
      return;
    }
    try {
      const results = await ImportProfilesFromURL(importUrl, importStrategy);
      toast.success(`Import complete: ${results.length} profiles processed.`);
      onClose();
    } catch (err) {
      toast.error(`Import failed: ${err}`);
    }
  };

  const toggleProfileSelection = (name: string) => {
    setSelectedProfiles((prev) =>
      prev.includes(name)
        ? prev.filter((p) => p !== name)
        : [...prev, name],
    );
  };

  if (!isOpen) {
    return null;
  }

  return (
    <div className="modal-overlay">
      <div className="modal-content">
        <div className="modal-header">
          <h2>Import/Export Profiles</h2>
          <button type="button" className="modal-close" onClick={onClose}>
            &times;
          </button>
        </div>
        <div className="modal-body">
          <div className="tab-buttons">
            <button
              type="button"
              className={`tab-button ${activeTab === "import" ? "active" : ""}`}
              onClick={() => setActiveTab("import")}
            >
              Import
            </button>
            <button
              type="button"
              className={`tab-button ${activeTab === "export" ? "active" : ""}`}
              onClick={() => setActiveTab("export")}
            >
              Export
            </button>
          </div>
          <div className="tab-content">
            {activeTab === "import" && (
              <div>
                <h3>Import from File or URL</h3>
                <div className="form-group">
                  <button type="button" className="btn-primary" onClick={handleImportFile}>
                    Import from File
                  </button>
                </div>
                <div className="form-group">
                  <label htmlFor="import-url">Or Enter URL</label>
                  <input
                    type="text"
                    id="import-url"
                    value={importUrl}
                    onChange={(e) => setImportUrl(e.target.value)}
                    placeholder="https://example.com/profiles.easyconfig"
                  />
                   <button type="button" className="btn-primary" onClick={handleImportUrl}>
                    Import from URL
                  </button>
                </div>
                <div className="form-group">
                  <h4>Conflict Resolution</h4>
                  <select
                    value={importStrategy}
                    onChange={(e) => setImportStrategy(e.target.value)}
                  >
                    <option value="skip">Skip Existing</option>
                    <option value="rename">Rename on Conflict</option>
                    <option value="overwrite">Overwrite Existing</option>
                  </select>
                </div>
              </div>
            )}
            {activeTab === "export" && (
              <div>
                <h3>Export Profiles</h3>
                <div className="profile-list">
                  {profiles.map((p) => (
                    <div key={p.name} className="profile-item">
                      <input
                        type="checkbox"
                        id={`profile-${p.name}`}
                        checked={selectedProfiles.includes(p.name)}
                        onChange={() => toggleProfileSelection(p.name)}
                      />
                      <label htmlFor={`profile-${p.name}`}>{p.name}</label>
                    </div>
                  ))}
                </div>
                <button
                  type="button"
                  className="btn-primary"
                  onClick={() => handleExport(false)}
                  disabled={selectedProfiles.length === 0}
                >
                  Export Selected
                </button>
                <button
                  type="button"
                  className="btn-secondary"
                  onClick={() => handleExport(true)}
                >
                  Export All
                </button>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default ImportExportModal;
