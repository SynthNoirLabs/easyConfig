import { act, fireEvent, render, screen } from "@testing-library/react";
import { beforeEach, expect, type Mock, test, vi } from "vitest";
import Marketplace from "../../../src/components/Marketplace";

const mockApp = {
  FetchPopularServers: vi.fn(),
  InstallMCPPackage: vi.fn(),
};

beforeEach(() => {
  window.go = {
    main: {
      App: mockApp,
    },
  } as Window["go"];
  mockApp.FetchPopularServers.mockReset();
  mockApp.InstallMCPPackage.mockReset();
});

const mockServers = [
  {
    name: "Server A",
    description: "Description A",
    source: JSON.stringify({ name: "Server A" }),
    verified: true,
    tags: ["tag1"],
  },
  {
    name: "Server B",
    description: "Description B",
    source: JSON.stringify({ name: "Server B" }),
    verified: false,
    tags: ["tag2"],
  },
];

test("Marketplace renders and filters servers", async () => {
  (window.go?.main?.App?.FetchPopularServers as Mock).mockResolvedValue(
    mockServers,
  );

  render(<Marketplace />);

  expect(await screen.findByText("Server A")).toBeInTheDocument();
  expect(await screen.findByText("Server B")).toBeInTheDocument();

  fireEvent.change(screen.getByPlaceholderText("Search packages..."), {
    target: { value: "Server A" },
  });
  expect(screen.queryByText("Server B")).not.toBeInTheDocument();
});

test("Marketplace handles installation", async () => {
  (window.go?.main?.App?.FetchPopularServers as Mock).mockResolvedValue(
    mockServers,
  );
  (window.go?.main?.App?.InstallMCPPackage as Mock).mockResolvedValue(
    undefined,
  );

  render(<Marketplace />);

  const installButtons = await screen.findAllByText("Install");
  await act(async () => {
    fireEvent.click(installButtons[0]);
  });

  expect(window.go?.main?.App?.InstallMCPPackage).toHaveBeenCalledWith(
    mockServers[0].source,
  );
});

test("Marketplace falls back to demo data without Wails", async () => {
  delete (window as Window & { go?: unknown }).go;

  render(<Marketplace />);

  expect(await screen.findByText("filesystem-mcp")).toBeInTheDocument();
  expect(
    screen.getByText("Discover and install extensions for your AI agents"),
  ).toBeInTheDocument();
});
