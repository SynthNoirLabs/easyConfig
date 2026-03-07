import { act, renderHook, waitFor } from "@testing-library/react";
import { beforeEach, expect, type Mock, test, vi } from "vitest";
import { ConfigProvider, useConfig } from "../../../src/context/ConfigContext";

const mockApp = {
  DiscoverConfigs: vi.fn(),
  ReadConfig: vi.fn(),
  SaveConfig: vi.fn(),
  DeleteConfig: vi.fn(),
};

beforeEach(() => {
  window.go = {
    main: {
      App: mockApp,
    },
  } as Window["go"];

  window.runtime = {
    EventsOn: vi.fn(() => () => {}),
    EventsOnMultiple: vi.fn(() => () => {}),
    EventsOnce: vi.fn(),
    EventsEmit: vi.fn(),
    WindowSetTitle: vi.fn(),
  } as Window["runtime"];

  mockApp.DiscoverConfigs.mockReset();
  mockApp.ReadConfig.mockReset();
  mockApp.SaveConfig.mockReset();
  mockApp.DeleteConfig.mockReset();
});

const mockConfigs = [
  {
    name: "test.json",
    path: "/path/to/test.json",
    format: "json",
    provider: "Test",
    fileName: "test.json",
    scope: "global",
  },
];

test("ConfigProvider fetches configs on mount", async () => {
  (window.go?.main?.App?.DiscoverConfigs as Mock).mockResolvedValue(
    mockConfigs,
  );

  const { result } = renderHook(() => useConfig(), {
    wrapper: ({ children }) => <ConfigProvider>{children}</ConfigProvider>,
  });

  await waitFor(() => {
    expect(result.current.configs).toEqual(mockConfigs);
    expect(result.current.loading).toBe(false);
  });
});

test("ConfigProvider handles deleting configs", async () => {
  (window.go?.main?.App?.DiscoverConfigs as Mock).mockResolvedValue(
    mockConfigs,
  );
  (window.go?.main?.App?.DeleteConfig as Mock).mockResolvedValue(undefined);

  const { result } = renderHook(() => useConfig(), {
    wrapper: ({ children }) => <ConfigProvider>{children}</ConfigProvider>,
  });

  // Wait for the initial fetch to complete
  await waitFor(() => {
    expect(result.current.configs.length).toBe(1);
  });
  expect(window.go?.main?.App?.DiscoverConfigs).toHaveBeenCalledTimes(1);

  // Trigger the delete action
  await act(async () => {
    await result.current.deleteConfig("/path/to/test.json");
  });

  expect(window.go?.main?.App?.DeleteConfig).toHaveBeenCalledWith(
    "/path/to/test.json",
  );
  // It refetches after delete
  expect(window.go?.main?.App?.DiscoverConfigs).toHaveBeenCalledTimes(2);
});

test("ConfigProvider falls back to demo data without Wails", async () => {
  delete (window as Window & { go?: unknown }).go;

  const { result } = renderHook(() => useConfig(), {
    wrapper: ({ children }) => <ConfigProvider>{children}</ConfigProvider>,
  });

  await waitFor(() => {
    expect(result.current.loading).toBe(false);
    expect(result.current.error).toBeNull();
    expect(result.current.configs[0]?.name).toBe("Desktop MCP Config");
  });

  await expect(
    result.current.readConfig("/demo/.claude/claude_desktop_config.json"),
  ).resolves.toContain("mcpServers");
});
