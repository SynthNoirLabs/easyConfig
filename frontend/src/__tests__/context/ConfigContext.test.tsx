import { act, renderHook, waitFor } from "@testing-library/react";
import React from "react";
import { expect, test, vi } from "vitest";
import { ConfigProvider, useConfig } from "../../../src/context/ConfigContext";

// Mock Wails Go backend
window.go = {
  main: {
    App: {
      DiscoverConfigs: vi.fn(),
      ReadConfig: vi.fn(),
      SaveConfig: vi.fn(),
      DeleteConfig: vi.fn(),
    },
  },
};

// Mock Wails JS runtime
window.runtime = {
  EventsOn: vi.fn(() => () => {}), // Mock EventsOn to return a dummy 'off' function
  EventsOnMultiple: vi.fn(() => () => {}),
  EventsOnce: vi.fn(),
  EventsEmit: vi.fn(),
  WindowSetTitle: vi.fn(),
  // Add other runtime functions if needed
};

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
  window.go.main.App.DiscoverConfigs.mockResolvedValue(mockConfigs);

  const { result } = renderHook(() => useConfig(), {
    wrapper: ({ children }) => <ConfigProvider>{children}</ConfigProvider>,
  });

  await waitFor(() => {
    expect(result.current.configs).toEqual(mockConfigs);
    expect(result.current.loading).toBe(false);
  });
});

test("ConfigProvider handles deleting configs", async () => {
  // Reset call counts for this test
  window.go.main.App.DiscoverConfigs.mockClear();
  window.go.main.App.DiscoverConfigs.mockResolvedValue(mockConfigs);
  window.go.main.App.DeleteConfig.mockResolvedValue(undefined);

  const { result } = renderHook(() => useConfig(), {
    wrapper: ({ children }) => <ConfigProvider>{children}</ConfigProvider>,
  });

  // Wait for the initial fetch to complete
  await waitFor(() => {
    expect(result.current.configs.length).toBe(1);
  });
  expect(window.go.main.App.DiscoverConfigs).toHaveBeenCalledTimes(1);

  // Trigger the delete action
  await act(async () => {
    await result.current.deleteConfig("/path/to/test.json");
  });

  expect(window.go.main.App.DeleteConfig).toHaveBeenCalledWith(
    "/path/to/test.json",
  );
  // It refetches after delete
  expect(window.go.main.App.DiscoverConfigs).toHaveBeenCalledTimes(2);
});
