import { useState } from "react";

type WizardStepId = "selectProvider" | "chooseScope" | "review";

interface WizardStep {
  id: WizardStepId;
  title: string;
  description: string;
}

const steps: WizardStep[] = [
  {
    id: "selectProvider",
    title: "Select a provider",
    description:
      "Choose which tool you want to configure (e.g. Claude Code, Gemini CLI, Copilot, Codex).",
  },
  {
    id: "chooseScope",
    title: "Choose configuration scope",
    description:
      "Decide whether this config should live at the user, project, or system level.",
  },
  {
    id: "review",
    title: "Review and create",
    description:
      "Review the summary of your choices before generating a configuration file.",
  },
];

const providers = [
  "Claude Code",
  "Gemini CLI",
  "GitHub Copilot",
  "Codex CLI",
  "Jules",
];

const scopes = ["User (home directory)", "Project (repo root)", "System"];

export default function ConfigWizard() {
  const [currentStepIndex, setCurrentStepIndex] = useState(0);
  const [selectedProvider, setSelectedProvider] = useState<string | null>(null);
  const [selectedScope, setSelectedScope] = useState<string | null>(null);

  const currentStep = steps[currentStepIndex];
  const isFirstStep = currentStepIndex === 0;
  const isLastStep = currentStepIndex === steps.length - 1;

  const canGoNext =
    (currentStep.id === "selectProvider" && selectedProvider !== null) ||
    (currentStep.id === "chooseScope" && selectedScope !== null) ||
    currentStep.id === "review";

  const handleNext = () => {
    if (!canGoNext || isLastStep) {
      return;
    }
    setCurrentStepIndex((index) => Math.min(index + 1, steps.length - 1));
  };

  const handleBack = () => {
    if (isFirstStep) {
      return;
    }
    setCurrentStepIndex((index) => Math.max(index - 1, 0));
  };

  return (
    <section className="config-wizard">
      <h3>Quick Config Wizard</h3>
      <p className="config-wizard__subtitle">
        Start a guided flow for one common configuration scenario. This is a
        non-destructive helper; you can cancel at any time.
      </p>

      <ol className="config-wizard__steps">
        {steps.map((step, index) => (
          <li
            key={step.id}
            className={`config-wizard__step${
              index === currentStepIndex ? " config-wizard__step--active" : ""
            }`}
          >
            <span className="config-wizard__step-index">{index + 1}</span>
            <div>
              <div className="config-wizard__step-title">{step.title}</div>
              <div className="config-wizard__step-description">
                {step.description}
              </div>
            </div>
          </li>
        ))}
      </ol>

      <div className="config-wizard__body">
        {currentStep.id === "selectProvider" && (
          <div className="config-wizard__choices">
            {providers.map((provider) => (
              <button
                key={provider}
                type="button"
                className={`config-wizard__choice${
                  selectedProvider === provider
                    ? " config-wizard__choice--selected"
                    : ""
                }`}
                onClick={() => setSelectedProvider(provider)}
              >
                {provider}
              </button>
            ))}
          </div>
        )}

        {currentStep.id === "chooseScope" && (
          <div className="config-wizard__choices">
            {scopes.map((scope) => (
              <button
                key={scope}
                type="button"
                className={`config-wizard__choice${
                  selectedScope === scope
                    ? " config-wizard__choice--selected"
                    : ""
                }`}
                onClick={() => setSelectedScope(scope)}
              >
                {scope}
              </button>
            ))}
          </div>
        )}

        {currentStep.id === "review" && (
          <div className="config-wizard__review">
            <p>
              <strong>Provider:</strong>{" "}
              {selectedProvider ?? "Not selected yet"}
            </p>
            <p>
              <strong>Scope:</strong> {selectedScope ?? "Not selected yet"}
            </p>
            <p className="config-wizard__hint">
              When you&apos;re ready, you can turn this summary into a real
              configuration by wiring it up to the backend and provider creation
              APIs.
            </p>
          </div>
        )}
      </div>

      <div className="config-wizard__footer">
        <button
          type="button"
          className="config-wizard__button"
          onClick={handleBack}
          disabled={isFirstStep}
        >
          Back
        </button>
        <button
          type="button"
          className="config-wizard__button config-wizard__button--primary"
          onClick={handleNext}
          disabled={!canGoNext || isLastStep}
        >
          Next
        </button>
      </div>
    </section>
  );
}
