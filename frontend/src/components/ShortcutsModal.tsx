import type React from "react";
import "./ShortcutsModal.css";

interface ShortcutsModalProps {
  isOpen: boolean;
  onClose: () => void;
}

const shortcuts = [
  {
    category: "Global",
    items: [
      { key: "Ctrl/Cmd + S", action: "Save current config" },
      { key: "Ctrl/Cmd + N", action: "Add new configuration" },
    ],
  },
  {
    category: "Navigation",
    items: [
      { key: "Ctrl/Cmd + 1", action: "Switch to Code view" },
      { key: "Ctrl/Cmd + 2", action: "Switch to Form view" },
      { key: "Ctrl/Cmd + 3", action: "Switch to Preview view" },
    ],
  },
  {
    category: "Editor",
    items: [
      { key: "Ctrl/Cmd + Z", action: "Undo" },
      { key: "Ctrl/Cmd + Shift + Z", action: "Redo" },
      { key: "Ctrl/Cmd + F", action: "Find in file" },
      { key: "Ctrl/Cmd + /", action: "Toggle comment" },
    ],
  },
];

const ShortcutsModal: React.FC<ShortcutsModalProps> = ({ isOpen, onClose }) => {
  if (!isOpen) {
    return null;
  }

  return (
    <div className="shortcuts-modal-overlay" onClick={onClose}>
      <div
        className="shortcuts-modal-content"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="shortcuts-modal-header">
          <h2>Keyboard Shortcuts</h2>
          <button type="button" className="close-button" onClick={onClose}>
            &times;
          </button>
        </div>
        <div className="shortcuts-modal-body">
          {shortcuts.map((category) => (
            <div key={category.category} className="shortcut-category">
              <h3>{category.category}</h3>
              <ul>
                {category.items.map((item) => (
                  <li key={item.key}>
                    <span className="shortcut-key">{item.key}</span>
                    <span className="shortcut-action">{item.action}</span>
                  </li>
                ))}
              </ul>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default ShortcutsModal;
