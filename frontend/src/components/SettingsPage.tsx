import { useEffect, useState } from "react";
import { GetSettings, SaveSettings } from "../../wailsjs/go/main/App";
import type { settings } from "../../wailsjs/go/models";

import "./SettingsPage.css";

export function SettingsPage() {
  const [scanDirs, setScanDirs] = useState<string[]>([]);
  const [newDir, setNewDir] = useState("");

  useEffect(() => {
    async function fetchSettings() {
      try {
        const s = await GetSettings();
        setScanDirs(s.providerScanDirs || []);
      } catch (error) {
        console.error("Failed to fetch settings:", error);
      }
    }
    fetchSettings();
  }, []);

  async function handleSave() {
    try {
      await SaveSettings({ providerScanDirs: scanDirs });
      alert("Settings saved successfully!");
    } catch (error) {
      console.error("Failed to save settings:", error);
      alert("Failed to save settings.");
    }
  }

  function handleAddDir() {
    if (newDir && !scanDirs.includes(newDir)) {
      setScanDirs([...scanDirs, newDir]);
      setNewDir("");
    }
  }

  function handleRemoveDir(dir: string) {
    setScanDirs(scanDirs.filter((d) => d !== dir));
  }

  return (
    <div className="settings-page">
      <h2>Settings</h2>
      <div className="settings-group">
        <h3>Dynamic Provider Discovery</h3>
        <p>Add directories to scan for dynamic provider definitions (provider.yaml).</p>
        <div className="scan-dirs-list">
          {scanDirs.map((dir) => (
            <div key={dir} className="scan-dir-item">
              <span>{dir}</span>
              <button onClick={() => handleRemoveDir(dir)}>Remove</button>
            </div>
          ))}
        </div>
        <div className="add-scan-dir">
          <input
            type="text"
            value={newDir}
            onChange={(e) => setNewDir(e.target.value)}
            placeholder="Enter directory path"
          />
          <button onClick={handleAddDir}>Add</button>
        </div>
      </div>
      <div className="settings-actions">
        <button onClick={handleSave}>Save Settings</button>
      </div>
    </div>
  );
}
