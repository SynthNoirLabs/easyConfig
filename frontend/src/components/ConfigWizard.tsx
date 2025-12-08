import React, { useState, useEffect } from 'react';
import './ConfigWizard.css';
import { StartWizard, NextWizardStep, CancelWizard } from '../../wailsjs/go/main/App';
import type { config } from '../../wailsjs/go/models';

interface ConfigWizardProps {
  providerName: string;
  isOpen: boolean;
  onClose: () => void;
}

const ConfigWizard: React.FC<ConfigWizardProps> = ({ providerName, isOpen, onClose }) => {
  const [currentStep, setCurrentStep] = useState<config.WizardStep | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [response, setResponse] = useState('');

  useEffect(() => {
    if (isOpen) {
      startWizard();
    }
  }, [isOpen]);

  const startWizard = async () => {
    try {
      const step = await StartWizard(providerName);
      setCurrentStep(step);
      setError(null);
    } catch (err) {
      setError(String(err));
    }
  };

  const handleNext = async () => {
    if (!currentStep) return;
    try {
      const nextStep = await NextWizardStep(providerName, currentStep.id, response);
      setCurrentStep(nextStep);
      setResponse('');
      setError(null);
    } catch (err) {
      setError(String(err));
    }
  };

  const handleCancel = async () => {
    try {
      await CancelWizard(providerName);
      onClose();
    } catch (err) {
      setError(String(err));
    }
  };

  if (!isOpen) {
    return null;
  }

  return (
    <div className="wizard-modal">
      <div className="wizard-content">
        {currentStep ? (
          <>
            <h2>{currentStep.title}</h2>
            <p>{currentStep.description}</p>
            {currentStep.id !== 'welcome' && currentStep.id !== 'done' && (
              <input
                type="text"
                value={response}
                onChange={(e) => setResponse(e.target.value)}
                placeholder="Enter your response"
              />
            )}
            <div className="wizard-buttons">
              {currentStep.id !== 'done' ? (
                <button onClick={handleNext}>Next</button>
              ) : (
                <button onClick={onClose}>Finish</button>
              )}
              <button onClick={handleCancel}>Cancel</button>
            </div>
          </>
        ) : (
          <p>Loading wizard...</p>
        )}
        {error && <p className="error">{error}</p>}
      </div>
    </div>
  );
};

export default ConfigWizard;
