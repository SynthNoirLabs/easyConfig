import { useState, useEffect } from "react";
import { toast } from "sonner";
import { GenerateWorkflow, SaveWorkflow, GetSupportedWorkflows, SetSecret } from "../../wailsjs/go/main/App";
import { Bot, Terminal, Code2, Github, Check, ArrowRight } from "lucide-react";
import "./Workflows.css";

export default function Workflows() {
  const [agents, setAgents] = useState<string[]>([]);
  const [selectedAgent, setSelectedAgent] = useState("");
  const [selectedTrigger, setSelectedTrigger] = useState("");
  const [generatedContent, setGeneratedContent] = useState("");
  const [filename, setFilename] = useState("");
  
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

  const handleGenerate = async (agent: string, trigger: string) => {
    try {
      const [content, secrets, instructions] = await GenerateWorkflow(agent, trigger);
      setGeneratedContent(content);
      setRequiredSecrets(secrets || []);
      setSetupInstructions(instructions || "");
      setSecretValues({}); 
      
      const agentLower = agent.toLowerCase();
      const triggerLower = trigger.toLowerCase();
      setFilename(`${agentLower}-${triggerLower}.yml`);
      
      toast.success("Workflow generated!");
    } catch (err) {
      toast.error("Failed to generate workflow: " + err);
    }
  };

  const handleSave = async () => {
    if (!filename || !generatedContent) return;
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

  // Icons mapping
  const getIcon = (name: string) => {
    if (name.includes("Claude")) return <Bot size={24} />;
    if (name.includes("Jules")) return <Terminal size={24} />;
    if (name.includes("Codex")) return <Code2 size={24} />;
    if (name.includes("Copilot")) return <Github size={24} />;
    return <Bot size={24} />;
  };

  return (
    <div className="workflows-container">
      <div className="workflows-header">
        <h2>Workflow Generator</h2>
        <p>Automate your AI agents with GitHub Actions CI/CD pipelines.</p>
      </div>

      {!generatedContent ? (
        <div className="agents-grid">
          {options.map((opt, i) => (
            <div 
              key={i} 
              className="agent-card"
              onClick={() => {
                setSelectedAgent(opt.agent);
                setSelectedTrigger(opt.trigger);
                handleGenerate(opt.agent, opt.trigger);
              }}
            >
              <div className="agent-icon">
                {getIcon(opt.agent)}
              </div>
              <div className="agent-info">
                <h3>{opt.agent}</h3>
                <span className="trigger-badge">{opt.trigger} Trigger</span>
              </div>
              <div className="agent-arrow">
                <ArrowRight size={20} />
              </div>
            </div>
          ))}
        </div>
      ) : (
        <div className="workflow-workspace">
          <button className="btn-back" onClick={() => setGeneratedContent("")}>
            ← Back to Agents
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
                <button className="btn btn-primary" onClick={handleSave}>
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
                  <p className="secrets-hint">Set these secrets in your repository.</p>
                  
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
                            className="input"
                          />
                          <button 
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
