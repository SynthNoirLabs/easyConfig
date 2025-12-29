import { fireEvent, render, screen } from "@testing-library/react";
import React from "react";
import { test, expect, vi } from "vitest";
import ConfigEditor from "../../../src/components/ConfigEditor";
import { useConfig } from "../../../src/context/ConfigContext";

// Mock hooks and components
vi.mock("../../../src/context/ConfigContext");
vi.mock("@monaco-editor/react", () => ({
  __esModule: true,
  default: ({ value, onChange }) => (
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

const mockReadConfig = vi.fn();
const mockSaveConfig = vi.fn();

const mockConfigItem = {
  name: "test.json",
  path: "/path/to/test.json",
  format: "json",
  provider: "Test",
  fileName: "test.json",
  scope: "global",
};

test("ConfigEditor renders and loads content", async () => {
  mockReadConfig.mockResolvedValue(`{"key": "value"}`);
  useConfig.mockReturnValue({ readConfig: mockReadConfig });

  render(<ConfigEditor configItem={mockConfigItem} />);

  expect(await screen.findByText("test.json")).toBeInTheDocument();
  expect(await screen.findByDisplayValue(`{"key": "value"}`)).toBeInTheDocument();
});

test("ConfigEditor handles editing and saving", async () => {
  mockReadConfig.mockResolvedValue(`{"key": "value"}`);
  mockSaveConfig.mockResolvedValue(undefined);
  useConfig.mockReturnValue({
    readConfig: mockReadConfig,
    saveConfig: mockSaveConfig,
  });

  render(<ConfigEditor configItem={mockConfigItem} />);
  const editor = await screen.findByTestId("monaco-editor");

  fireEvent.change(editor, { target: { value: `{"key": "new_value"}` } });
  expect(screen.getByText("Save")).not.toBeDisabled();
  fireEvent.click(screen.getByText("Save"));

  expect(mockSaveConfig).toHaveBeenCalledWith(
    "/path/to/test.json",
    `{"key": "new_value"}`,
  );
});
