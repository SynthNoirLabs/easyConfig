import { useState, useEffect } from "react";
import { toast } from "sonner";
import { GenerateWorkflow, SaveWorkflow, GetSupportedWorkflows } from "../../wailsjs/go/main/App";
import "./Workflows.css";

export default function Workflows() {
  const [agents, setAgents] = useState<string[]>([]);
  const [selectedAgent, setSelectedAgent] = useState("");
  const [selectedTrigger, setSelectedTrigger] = useState("");
  const [generatedContent, setGeneratedContent] = useState("");
  const [filename, setFilename] = useState("");

  useEffect(() => {
    loadSupportedWorkflows();
  }, []);

  const loadSupportedWorkflows = async () => {
    try {
      const workflows = await GetSupportedWorkflows();
      setAgents(workflows);
      if (workflows.length > 0) {
        // Parse "Agent (Trigger)" format
        const first = workflows[0];
        const parts = first.split(" (");
        if (parts.length > 1) {
             // Just setting defaults for UI
        }
      }
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
      const content = await GenerateWorkflow(selectedAgent, selectedTrigger);
      setGeneratedContent(content);
      
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

  // Helper to parse the "Agent (Trigger)" strings from backend
  // In a real app, backend might return structured objects
  const parseOptions = () => {
      // For this MVP, we'll hardcode the mapping based on what backend returns
      // Backend returns strings like "Claude (Comment)"
      // We can just let user select from the list
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
                if (idx > 0) { // 0 is placeholder
                    const opt = options[idx - 1];
                    setSelectedAgent(opt.agent);
                    setSelectedTrigger(opt.trigger);
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
      )}
    </div>
  );
}
