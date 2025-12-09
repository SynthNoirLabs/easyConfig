import type React from "react";
import { useEffect, useMemo, useState } from "react";
import { CreateConfig } from "../../wailsjs/go/main/App";
import { useConfig } from "../context/ConfigContext";
import "./QuickConfigWizard.css";

interface QuickConfigWizardProps {
  isOpen: boolean;
  onClose: () => void;
  onSuccess: () => void;
}

const DEFAULT_PROVIDERS = [
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

const QuickConfigWizard: React.FC<QuickConfigWizardProps> = ({
  isOpen,
  onClose,
  onSuccess,
}) => {
  const { configs } = useConfig();
  const providerOptions = useMemo(() => {
    const discovered = Array.from(new Set(configs.map((c) => c.provider)));
    const combined = [...DEFAULT_PROVIDERS, ...discovered];
    return Array.from(new Set(combined));
  }, [configs]);

  const [step, setStep] = useState(1);
  const [provider, setProvider] = useState(
    providerOptions[0] ?? DEFAULT_PROVIDERS[0],
  );
  const [scope, setScope] = useState(SCOPES[0].value);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!providerOptions.includes(provider)) {
      setProvider(providerOptions[0] ?? DEFAULT_PROVIDERS[0]);
    }
  }, [providerOptions, provider]);

  if (!isOpen) return null;

  const handleNext = () => {
    setStep(step + 1);
  };

  const handleBack = () => {
    setStep(step - 1);
  };

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
          <h3>Quick Config Wizard</h3>
          <button type="button" className="btn-close" onClick={onClose}>
            &times;
          </button>
        </div>
        <form onSubmit={handleSubmit}>
          <div className="modal-body">
            {error && <div className="modal-error">{error}</div>}

            {step === 1 && (
              <div className="form-group">
                <label htmlFor="provider">1. Choose a provider</label>
                <p className="form-group-description">
                  Select the tool you want to configure.
                </p>
                <select
                  id="provider"
                  value={provider}
                  onChange={(e) => setProvider(e.target.value)}
                >
                  {providerOptions.map((p) => (
                    <option key={p} value={p}>
                      {p}
                    </option>
                  ))}
                </select>
              </div>
            )}

            {step === 2 && (
              <div className="form-group">
                <label htmlFor="scope">2. Set the scope</label>
                <p className="form-group-description">
                  Global applies to your user account. Project applies only to
                  the current workspace.
                </p>
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
            )}
          </div>
          <div className="modal-footer">
            {step === 1 && (
              <>
                <button
                  type="button"
                  className="btn-cancel"
                  onClick={onClose}
                  disabled={isSubmitting}
                >
                  Cancel
                </button>
                <button
                  type="button"
                  className="btn-submit"
                  onClick={handleNext}
                >
                  Next
                </button>
              </>
            )}
            {step === 2 && (
              <>
                <button
                  type="button"
                  className="btn-cancel"
                  onClick={handleBack}
                >
                  Back
                </button>
                <button
                  type="submit"
                  className="btn-submit"
                  disabled={isSubmitting}
                >
                  {isSubmitting ? "Creating..." : "Create"}
                </button>
              </>
            )}
          </div>
        </form>
      </div>
    </div>
  );
};

export default QuickConfigWizard;
