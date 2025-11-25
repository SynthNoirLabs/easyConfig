import { useState, useEffect } from "react";
import { toast } from "sonner";
import { GenerateWorkflow, SaveWorkflow, GetSupportedWorkflows, SetSecret } from "../../wailsjs/go/main/App";
import "./Workflows.css";

export default function Workflows() {
  const [agents, setAgents] = useState<string[]>([]);
  const [selectedAgent, setSelectedAgent] = useState("");
  const [selectedTrigger, setSelectedTrigger] = useState("");
  const [generatedContent, setGeneratedContent] = useState("");
  const [filename, setFilename] = useState("");
  
  // New state for secrets and nudges
  const [requiredSecrets, setRequiredSecrets] = useState<string[]>([]);
  const [setupInstructions, setSetupInstructions] = useState("");
  const [secretValues, setSecretValues] = useState<Record<string, string>>({});

  useEffect(() => {
    loadSupportedWorkflows();
  }, []);

  const loadSupportedWorkflows = async () => {
    try {
      const workflows = await GetSupportedWorkflows();
      setAgents(workflows);
    } catch (err) {
      toast.error("Failed to load supported workflows");
    }
  };

  const handleGenerate = async () => {
    if (!selectedAgent || !selectedTrigger) {
      toast.error("Please select an agent and a trigger");
      return;
    }

    try {
      // Backend now returns: content, requiredSecrets, setupInstructions, error
      const [content, secrets, instructions] = await GenerateWorkflow(selectedAgent, selectedTrigger);
      setGeneratedContent(content);
      setRequiredSecrets(secrets || []);
      setSetupInstructions(instructions || "");
      setSecretValues({}); // Reset secret inputs
      
      // Auto-generate filename
      const agentLower = selectedAgent.toLowerCase();
      const triggerLower = selectedTrigger.toLowerCase();
      setFilename(`${agentLower}-${triggerLower}.yml`);
      
      toast.success("Workflow generated!");
    } catch (err) {
      toast.error("Failed to generate workflow: " + err);
    }
  };

  const handleSave = async () => {
    if (!filename || !generatedContent) {
      return;
    }

    try {
      await SaveWorkflow(filename, generatedContent);
      toast.success(`Saved to .github/workflows/${filename}`);
    } catch (err) {
      toast.error("Failed to save workflow: " + err);
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
      // Clear the input for security
      setSecretValues(prev => ({ ...prev, [secretName]: "" }));
    } catch (err) {
      toast.error(`Failed to set secret: ${err}`);
    }
  };

  const parseOptions = () => {
      return agents.map(a => {
          const match = a.match(/(.+) \((.+)\)/);
          if (match) {
              return { label: a, agent: match[1], trigger: match[2] };
          }
          return { label: a, agent: a, trigger: "Manual" };
      });
  };

  const options = parseOptions();

  return (
    <div className="workflows-container">
      <h2>GitHub Actions Generator</h2>
      <p>Create CI/CD workflows for your AI agents.</p>

      <div className="controls">
        <div className="control-group">
          <label>Select Workflow Type</label>
          <select 
            onChange={(e) => {
                const idx = e.target.selectedIndex;
                if (idx > 0) {
                    const opt = options[idx - 1];
                    setSelectedAgent(opt.agent);
                    setSelectedTrigger(opt.trigger);
                    // Reset state when selection changes
                    setGeneratedContent("");
                    setRequiredSecrets([]);
                    setSetupInstructions("");
                }
            }}
          >
            <option value="">-- Select --</option>
            {options.map((opt, i) => (
              <option key={i} value={opt.label}>{opt.label}</option>
            ))}
          </select>
        </div>

        <button className="btn-primary" onClick={handleGenerate} disabled={!selectedAgent}>
          Generate
        </button>
      </div>

      {generatedContent && (
        <div className="workflow-content">
          <div className="preview-section">
            <div className="preview-header">
              <input 
                type="text" 
                value={filename} 
                onChange={(e) => setFilename(e.target.value)}
                placeholder="workflow.yml"
              />
              <button className="btn-secondary" onClick={handleSave}>
                Save to Project
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
                <p className="secrets-hint">Set these secrets in your repository to enable the workflow.</p>
                
                <div className="secrets-list">
                  {requiredSecrets.map(secret => (
                    <div key={secret} className="secret-item">
                      <label>{secret}</label>
                      <div className="secret-input-group">
                        <input 
                          type="password" 
                          placeholder="Enter value..."
                          value={secretValues[secret] || ""}
                          onChange={(e) => setSecretValues(prev => ({ ...prev, [secret]: e.target.value }))}
                        />
                        <button 
                          className="btn-small"
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
      )}
    </div>
  );
}
