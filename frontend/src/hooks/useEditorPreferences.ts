import { useState } from "react";

export interface EditorPreferences {
  fontSize: number;
  fontFamily: string;
  wordWrap: "off" | "on" | "wordWrapColumn";
  minimap: boolean;
  lineNumbers: "on" | "off" | "relative";
  tabSize: number;
  insertSpaces: boolean;
  theme: "vs-dark" | "light" | "high-contrast";
  renderWhitespace: "none" | "boundary" | "all";
}

export const defaultPreferences: EditorPreferences = {
  fontSize: 14,
  fontFamily: "Fira Code",
  wordWrap: "off",
  minimap: false,
  lineNumbers: "on",
  tabSize: 2,
  insertSpaces: true,
  theme: "vs-dark",
  renderWhitespace: "none",
};

export function useEditorPreferences(): [
  EditorPreferences,
  (updates: Partial<EditorPreferences>) => void,
] {
  const [prefs, setPrefs] = useState<EditorPreferences>(() => {
    try {
      const saved = localStorage.getItem("editorPreferences");
      return saved
        ? { ...defaultPreferences, ...JSON.parse(saved) }
        : defaultPreferences;
    } catch (error) {
      console.error("Failed to parse editor preferences:", error);
      return defaultPreferences;
    }
  });

  const updatePrefs = (updates: Partial<EditorPreferences>) => {
    setPrefs((prev) => {
      const newPrefs = { ...prev, ...updates };
      try {
        localStorage.setItem("editorPreferences", JSON.stringify(newPrefs));
      } catch (error) {
        console.error("Failed to save editor preferences:", error);
      }
      return newPrefs;
    });
  };

  return [prefs, updatePrefs];
}
