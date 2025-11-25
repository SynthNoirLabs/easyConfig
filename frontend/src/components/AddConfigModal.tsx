import type React from "react";
import { useState } from "react";
import { CreateConfig } from "../../wailsjs/go/main/App";
import "./AddConfigModal.css";

interface AddConfigModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSuccess: () => void;
}

const PROVIDERS = [
  "Claude Code",
  "Gemini",
  "OpenAI",
  "Codex CLI",
  "GitHub Copilot",
  "OpenCode",
  "Crush CLI",
  "Jules",
  "Git",
  "Aider",
  "Goose",
];

const SCOPES = [
  { value: "global", label: "Global (User Home)" },
  { value: "project", label: "Project (Current Workspace)" },
];

const AddConfigModal: React.FC<AddConfigModalProps> = ({
  isOpen,
  onClose,
  onSuccess,
}) => {
  const [provider, setProvider] = useState(PROVIDERS[0]);
  const [scope, setScope] = useState(SCOPES[0].value);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  if (!isOpen) return null;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);
    setError(null);

    try {
      // projectPath is managed by backend context for now, passing empty for global
      // For project scope, backend uses the opened project path
      await CreateConfig(provider, scope, "");
      onSuccess();
      onClose();
    } catch (err) {
      console.error(err);
      setError(err instanceof Error ? err.message : "Failed to create config");
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="modal-overlay">
      <div className="modal-content">
        <div className="modal-header">
          <h3>Add Configuration</h3>
          <button className="btn-close" onClick={onClose}>
            &times;
          </button>
        </div>
        <form onSubmit={handleSubmit}>
          <div className="modal-body">
            {error && <div className="modal-error">{error}</div>}

            <div className="form-group">
              <label htmlFor="provider">Provider</label>
              <select
                id="provider"
                value={provider}
                onChange={(e) => setProvider(e.target.value)}
              >
                {PROVIDERS.map((p) => (
                  <option key={p} value={p}>
                    {p}
                  </option>
                ))}
              </select>
            </div>

            <div className="form-group">
              <label htmlFor="scope">Scope</label>
              <select
                id="scope"
                value={scope}
                onChange={(e) => setScope(e.target.value)}
              >
                {SCOPES.map((s) => (
                  <option key={s.value} value={s.value}>
                    {s.label}
                  </option>
                ))}
              </select>
            </div>
          </div>
          <div className="modal-footer">
            <button
              type="button"
              className="btn-cancel"
              onClick={onClose}
              disabled={isSubmitting}
            >
              Cancel
            </button>
            <button
              type="submit"
              className="btn-submit"
              disabled={isSubmitting}
            >
              {isSubmitting ? "Creating..." : "Create"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default AddConfigModal;
