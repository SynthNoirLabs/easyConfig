import { fireEvent, render, screen } from "@testing-library/react";
import React from "react";
import { expect, test, vi } from "vitest";
import AddConfigModal from "../../../src/components/AddConfigModal";
import { useConfig } from "../../../src/context/ConfigContext";

// Mock hooks and Wails bridge
vi.mock("../../../src/context/ConfigContext");
window.go = {
  main: {
    App: {
      CreateConfig: vi.fn(),
    },
  },
};

const mockOnClose = vi.fn();
const mockOnSuccess = vi.fn();

test("AddConfigModal renders when open", () => {
  useConfig.mockReturnValue({ configs: [] });
  render(
    <AddConfigModal
      isOpen={true}
      onClose={mockOnClose}
      onSuccess={mockOnSuccess}
    />,
  );
  expect(screen.getByText("Add Configuration")).toBeInTheDocument();
});

test("AddConfigModal handles form submission", async () => {
  useConfig.mockReturnValue({ configs: [] });
  window.go.main.App.CreateConfig.mockResolvedValue(undefined);

  render(
    <AddConfigModal
      isOpen={true}
      onClose={mockOnClose}
      onSuccess={mockOnSuccess}
    />,
  );

  fireEvent.click(screen.getByText("Create"));

  expect(window.go.main.App.CreateConfig).toHaveBeenCalledWith(
    "Claude Code",
    "global",
    "",
  );
});
