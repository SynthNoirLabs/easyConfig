import {
  ArrowRight,
  Bot,
  Check,
  Code2,
  Github,
  Search,
  Terminal,
} from "lucide-react";
import { useCallback, useEffect, useMemo, useState } from "react";
import { toast } from "sonner";
import {
  ListWorkflowTemplates,
  SaveWorkflow,
  SetSecret,
} from "../../wailsjs/go/main/App";
import type { workflows } from "../../wailsjs/go/models";
import "./Workflows.css";

export default function Workflows() {
  const [templates, setTemplates] = useState<workflows.Template[]>([]);
  const [filter, setFilter] = useState("");
  const [selected, setSelected] = useState<workflows.Template | null>(null);
  const [filename, setFilename] = useState("");
  const [secretValues, setSecretValues] = useState<Record<string, string>>({});

  const requiredSecrets = selected?.requiredSecrets ?? [];
  const setupInstructions = selected?.setupInstructions ?? "";
  const generatedContent = selected?.content ?? "";

  const loadTemplates = useCallback(async () => {
    try {
      const data = await ListWorkflowTemplates();
      setTemplates(data);
    } catch (_err) {
      toast.error("Failed to load workflow templates");
    }
  }, []);

  useEffect(() => {
    loadTemplates();
  }, [loadTemplates]);

  const filtered = useMemo(() => {
    const q = filter.trim().toLowerCase();
    if (!q) return templates;
    return templates.filter((t) =>
      [t.name, t.description, t.agent, t.trigger, ...(t.tags || [])]
        .filter(Boolean)
        .some((v) => v.toLowerCase().includes(q)),
    );
  }, [filter, templates]);

  const handleSelect = (tmpl: workflows.Template) => {
    setSelected(tmpl);
    setFilename(tmpl.defaultFilename || `${tmpl.id}.yml`);
    setSecretValues({});
  };

  const handleSave = async () => {
    if (!selected) return;
    if (!filename || !generatedContent) return;
    try {
      await SaveWorkflow(filename, generatedContent);
      toast.success(`Saved to .github/workflows/${filename}`);
    } catch (err) {
      toast.error(`Failed to save workflow: ${err}`);
    }
  };

  const handleSetSecret = async (secretName: string) => {
    const value = secretValues[secretName];
    if (!value) {
      toast.error(`Please enter a value for ${secretName}`);
      return;
    }
    try {
      await SetSecret(secretName, value);
      toast.success(`Secret ${secretName} set successfully!`);
      setSecretValues((prev) => ({ ...prev, [secretName]: "" }));
    } catch (err) {
      toast.error(`Failed to set secret: ${err}`);
    }
  };

  const getIcon = (name: string) => {
    if (name.toLowerCase().includes("claude")) return <Bot size={24} />;
    if (name.toLowerCase().includes("jules")) return <Terminal size={24} />;
    if (name.toLowerCase().includes("codex")) return <Code2 size={24} />;
    if (name.toLowerCase().includes("copilot")) return <Github size={24} />;
    return <Bot size={24} />;
  };

  return (
    <div className="workflows-container">
      <div className="workflows-header">
        <div>
          <h2>Workflow Gallery</h2>
          <p>Pick a template, preview, and drop it into .github/workflows.</p>
        </div>
        <div className="search-box">
          <Search size={18} className="search-icon" />
          <input
            type="text"
            placeholder="Search by agent, trigger, or tag"
            value={filter}
            onChange={(e) => setFilter(e.target.value)}
            className="search-input"
          />
        </div>
      </div>

      {!selected ? (
        <div className="agents-grid">
          {filtered.map((tmpl) => (
            <button
              type="button"
              key={tmpl.id}
              className="agent-card"
              onClick={() => handleSelect(tmpl)}
              aria-label={`Select ${tmpl.name}`}
            >
              <div className="agent-icon">{getIcon(tmpl.agent)}</div>
              <div className="agent-info">
                <h3>{tmpl.name}</h3>
                <span className="trigger-badge">{tmpl.trigger} trigger</span>
                <p className="agent-desc">{tmpl.description}</p>
                <div className="agent-tags">
                  {tmpl.tags?.slice(0, 3).map((tag) => (
                    <span key={tag} className="tag-chip">
                      {tag}
                    </span>
                  ))}
                </div>
              </div>
              <div className="agent-arrow">
                <ArrowRight size={20} />
              </div>
            </button>
          ))}
          {filtered.length === 0 && (
            <div className="empty-state">
              <h3>No templates match that search</h3>
              <p>Try a different keyword or clear the filter.</p>
            </div>
          )}
        </div>
      ) : (
        <div className="workflow-workspace">
          <button
            type="button"
            className="btn-back"
            onClick={() => setSelected(null)}
          >
            ← Back to Gallery
          </button>

          <div className="workflow-content">
            <div className="preview-section">
              <div className="preview-header">
                <input
                  type="text"
                  value={filename}
                  onChange={(e) => setFilename(e.target.value)}
                  className="input"
                />
                <button
                  type="button"
                  className="btn btn-primary"
                  onClick={handleSave}
                >
                  <Check size={16} /> Save to Project
                </button>
              </div>
              <textarea
                className="code-preview"
                value={generatedContent}
                readOnly
              />
            </div>

            <div className="sidebar-section">
              {setupInstructions && (
                <div className="nudge-box">
                  <h4>Setup Instructions</h4>
                  <p>{setupInstructions}</p>
                </div>
              )}

              {requiredSecrets.length > 0 && (
                <div className="secrets-box">
                  <h4>Required Secrets</h4>
                  <p className="secrets-hint">
                    Set these secrets in your repository.
                  </p>

                  <div className="secrets-list">
                    {requiredSecrets.map((secret) => (
                      <div key={secret} className="secret-item">
                        <label
                          htmlFor={`secret-${secret.replace(/[^a-zA-Z0-9_-]/g, "-")}`}
                        >
                          {secret}
                        </label>
                        <div className="secret-input-group">
                          <input
                            id={`secret-${secret.replace(/[^a-zA-Z0-9_-]/g, "-")}`}
                            type="password"
                            placeholder="Enter value..."
                            value={secretValues[secret] || ""}
                            onChange={(e) =>
                              setSecretValues((prev) => ({
                                ...prev,
                                [secret]: e.target.value,
                              }))
                            }
                            className="input"
                          />
                          <button
                            type="button"
                            className="btn btn-primary btn-sm"
                            onClick={() => handleSetSecret(secret)}
                          >
                            Set
                          </button>
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
