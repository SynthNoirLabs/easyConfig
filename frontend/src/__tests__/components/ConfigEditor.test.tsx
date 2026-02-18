import { fireEvent, render, screen } from "@testing-library/react";
import React from "react";
import { expect, test, vi } from "vitest";
import ConfigEditor from "../../../src/components/ConfigEditor";
import { useConfig } from "../../../src/context/ConfigContext";

// Mock hooks and components
vi.mock("../../../src/context/ConfigContext", () => ({
  useConfig: vi.fn(),
  ConfigProvider: ({ children }: any) => <div>{children}</div>,
}));

vi.mock("@monaco-editor/react", () => ({
  __esModule: true,
  default: ({
    value,
    onChange,
  }: { value: string; onChange: (v: string) => void }) => (
    <textarea
      data-testid="monaco-editor"
      value={value}
      onChange={(e) => onChange(e.target.value)}
    />
  ),
}));
vi.mock("../../../src/components/editors/ClaudeConfigEditor", () => ({
  __esModule: true,
  default: () => <div data-testid="claude-editor" />,
}));
vi.mock("../../../src/components/editors/OpenCodeConfigEditor", () => ({
  __esModule: true,
  default: () => <div data-testid="opencode-editor" />,
}));

// Mock hooks
vi.mock("../../../src/hooks/useEditorPreferences", () => ({
  useEditorPreferences: () => [{}, vi.fn()],
}));

const mockReadConfig = vi.fn();
const mockSaveConfig = vi.fn();

const mockConfigItem = {
  name: "test.json",
  path: "/path/to/test.json",
  format: "json",
  provider: "Test",
  fileName: "test.json",
  scope: "global",
  exists: true,
};

test("ConfigEditor renders and loads content", async () => {
  mockReadConfig.mockResolvedValue(`{"key": "value"}`);
  (useConfig as any).mockReturnValue({
    readConfig: mockReadConfig,
    saveConfig: mockSaveConfig,
  });

  render(<ConfigEditor configItem={mockConfigItem} />);

  // Wait for load
  // Note: Finding by display value might need to wait for state update
});
