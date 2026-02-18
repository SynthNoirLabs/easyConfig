import type React from "react";
import type { EditorPreferences } from "../hooks/useEditorPreferences";
import "./EditorSettings.css";

interface EditorSettingsProps {
  preferences: EditorPreferences;
  onChange: (updates: Partial<EditorPreferences>) => void;
  onReset: () => void;
}

const EditorSettings: React.FC<EditorSettingsProps> = ({
  preferences,
  onChange,
  onReset,
}) => {
  const handleValueChange = (key: keyof EditorPreferences, value: unknown) => {
    onChange({ [key]: value });
  };

  return (
    <div className="editor-settings-dropdown">
      <div className="settings-group">
        <label htmlFor="font-size">Font Size</label>
        <div className="settings-control">
          <input
            type="range"
            id="font-size"
            min="10"
            max="24"
            value={preferences.fontSize}
            onChange={(e) =>
              handleValueChange("fontSize", Number(e.target.value))
            }
          />
          <span>{preferences.fontSize}px</span>
        </div>
      </div>

      <div className="settings-group">
        <label htmlFor="tab-size">Tab Size</label>
        <div className="settings-control">
          <select
            id="tab-size"
            value={preferences.tabSize}
            onChange={(e) =>
              handleValueChange("tabSize", Number(e.target.value))
            }
          >
            <option value={2}>2 Spaces</option>
            <option value={4}>4 Spaces</option>
          </select>
        </div>
      </div>

      <div className="settings-group">
        <label htmlFor="theme">Theme</label>
        <div className="settings-control">
          <select
            id="theme"
            value={preferences.theme}
            onChange={(e) =>
              handleValueChange(
                "theme",
                e.target.value as EditorPreferences["theme"],
              )
            }
          >
            <option value="vs-dark">Dark</option>
            <option value="light">Light</option>
            <option value="high-contrast">High Contrast</option>
          </select>
        </div>
      </div>

      <div className="settings-group">
        <label>
          <input
            type="checkbox"
            checked={preferences.minimap}
            onChange={(e) => handleValueChange("minimap", e.target.checked)}
          />
          Show Minimap
        </label>
      </div>

      <div className="settings-group">
        <label>
          <input
            type="checkbox"
            checked={preferences.wordWrap === "on"}
            onChange={(e) =>
              handleValueChange("wordWrap", e.target.checked ? "on" : "off")
            }
          />
          Word Wrap
        </label>
      </div>

      <div className="settings-group">
        <label>
          <input
            type="checkbox"
            checked={preferences.lineNumbers !== "off"}
            onChange={(e) =>
              handleValueChange("lineNumbers", e.target.checked ? "on" : "off")
            }
          />
          Line Numbers
        </label>
      </div>

      <div className="settings-footer">
        <button type="button" onClick={onReset} className="btn-secondary">
          Reset to Defaults
        </button>
      </div>
    </div>
  );
};

export default EditorSettings;
