import { X } from "lucide-react";
import type React from "react";
// import {
//   ExportAllProfiles,
//   ExportProfiles,
//   ImportProfilesFromFile,
//   ImportProfilesFromURL,
//   SaveExportedProfiles,
// } from "../../wailsjs/go/main/App";
import "./ImportExportModal.css";

interface ImportExportModalProps {
  isOpen: boolean;
  onClose: () => void;
}

const ImportExportModal: React.FC<ImportExportModalProps> = ({
  isOpen,
  onClose,
}) => {
  if (!isOpen) return null;

  return (
    <div className="modal-overlay">
      <div className="modal-content">
        <div className="modal-header">
          <h2>Import / Export Profiles</h2>
          <button type="button" className="close-button" onClick={onClose}>
            <X size={20} />
          </button>
        </div>
        <div className="modal-body">
          <p>Import/Export functionality is currently disabled.</p>
        </div>
        <div className="modal-footer">
          <button type="button" className="btn-secondary" onClick={onClose}>
            Close
          </button>
        </div>
      </div>
    </div>
  );
};

export default ImportExportModal;
