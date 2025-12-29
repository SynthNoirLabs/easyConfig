import { render, screen } from "@testing-library/react";
import React from "react";
import { test, expect, vi } from "vitest";
import Sidebar from "../../../src/components/Sidebar";

// Mock the useConfig hook
vi.mock("../../../src/context/ConfigContext", () => ({
  useConfig: () => ({
    deleteConfig: vi.fn(),
    readConfig: vi.fn(),
  }),
}));

// Mock Wails runtime
window.go = {
  main: {
    App: {
      ListProfiles: vi.fn().mockResolvedValue([]),
      GetProviderStatuses: vi.fn().mockResolvedValue([]),
    },
  },
};

// Mock ProviderStatusWidget to avoid useEffect issues in test
vi.mock("../../../src/components/ProviderStatusWidget", () => ({
  default: () => <div>ProviderStatusWidget</div>,
}));

test("Sidebar renders navigation buttons", () => {
  render(
    <Sidebar
      items={[]}
      onSelect={() => {}}
      onAdd={() => {}}
      currentView="configs"
      onViewChange={() => {}}
    />,
  );

  expect(screen.getByText("Dashboard")).toBeInTheDocument();
  expect(screen.getByText("Workflows")).toBeInTheDocument();
  expect(screen.getByText("Marketplace")).toBeInTheDocument();
  expect(screen.getByText("Docs")).toBeInTheDocument();
});
