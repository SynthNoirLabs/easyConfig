import { fireEvent, render, screen } from "@testing-library/react";
import React from "react";
import { expect, test, vi } from "vitest";
import Marketplace from "../../../src/components/Marketplace";

// Mock Wails bridge
window.go = {
  main: {
    App: {
      FetchPopularServers: vi.fn(),
      InstallMCPPackage: vi.fn(),
    },
  },
};

const mockServers = [
  {
    name: "Server A",
    description: "Description A",
    verified: true,
    tags: ["tag1"],
  },
  {
    name: "Server B",
    description: "Description B",
    verified: false,
    tags: ["tag2"],
  },
];

test("Marketplace renders and filters servers", async () => {
  window.go.main.App.FetchPopularServers.mockResolvedValue(mockServers);

  render(<Marketplace />);

  expect(await screen.findByText("Server A")).toBeInTheDocument();
  expect(await screen.findByText("Server B")).toBeInTheDocument();

  fireEvent.click(screen.getByText("Verified only"));
  expect(screen.queryByText("Server B")).not.toBeInTheDocument();
});

test("Marketplace handles installation", async () => {
  window.go.main.App.FetchPopularServers.mockResolvedValue(mockServers);
  window.go.main.App.InstallMCPPackage.mockResolvedValue(undefined);

  render(<Marketplace />);

  const installButtons = await screen.findAllByText("Install");
  fireEvent.click(installButtons[0]);

  expect(window.go.main.App.InstallMCPPackage).toHaveBeenCalledWith(
    JSON.stringify(mockServers[0]),
  );
});
