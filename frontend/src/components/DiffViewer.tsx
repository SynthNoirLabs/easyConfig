import { DiffEditor } from "@monaco-editor/react";
import type React from "react";

interface DiffViewerProps {
  original: string;
  modified: string;
  language: string;
}

const DiffViewer: React.FC<DiffViewerProps> = ({
  original,
  modified,
  language,
}) => {
  return (
    <DiffEditor
      height="100%"
      original={original}
      modified={modified}
      language={language}
      theme="vs-dark"
      options={{
        readOnly: true,
        renderSideBySide: true,
        minimap: { enabled: false },
        scrollBeyondLastLine: false,
        fontSize: 14,
        automaticLayout: true,
      }}
    />
  );
};

export default DiffViewer;
