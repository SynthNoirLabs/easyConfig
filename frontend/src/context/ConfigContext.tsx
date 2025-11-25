import type React from "react";
import { createContext, useCallback, useContext, useEffect, useState } from "react";
import { config } from "../../wailsjs/go/models";
import {
  DiscoverConfigs,
  ReadConfig,
  SaveConfig,
  CreateConfig,
  DeleteConfig,
} from "../../wailsjs/go/main/App";

interface ConfigContextType {
  configs: config.Item[];
  loading: boolean;
  error: string | null;
  refreshConfigs: () => Promise<void>;
  readConfig: (path: string) => Promise<string>;
  saveConfig: (path: string, content: string) => Promise<void>;
  deleteConfig: (path: string) => Promise<void>;
}

const ConfigContext = createContext<ConfigContextType | undefined>(undefined);

export const ConfigProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [configs, setConfigs] = useState<config.Item[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  const fetchConfigs = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      // DiscoverConfigs takes a projectPath. Empty string means discover in default locations.
      const items = await DiscoverConfigs("");
      setConfigs(items || []);
    } catch (err) {
      console.error("Failed to load configs:", err);
      setError(
        err instanceof Error ? err.message : "Failed to load configurations",
      );
      // Mock data for development if Wails is not available (e.g. in browser)
      if (
        String(err).includes("window.go") ||
        String(err).includes("is not a function")
      ) {
        console.warn("Wails runtime not found. Using mock data.");
        // Optional: Set mock data here if we want to test UI without backend
      }
    } finally {
      setLoading(false);
    }
  }, []);

  const readConfig = async (path: string): Promise<string> => {
    try {
      return await ReadConfig(path);
    } catch (err) {
      console.error("Failed to read config:", err);
      throw err;
    }
  };

  const saveConfig = async (path: string, content: string): Promise<void> => {
    try {
      await SaveConfig(path, content);
    } catch (err) {
      console.error("Failed to save config:", err);
      throw err;
    }
  };

  const deleteConfig = async (path: string): Promise<void> => {
    try {
      await DeleteConfig(path);
      await fetchConfigs(); // Refresh list after delete
    } catch (err) {
      console.error("Failed to delete config:", err);
      throw err;
    }
  };

  useEffect(() => {
    fetchConfigs();
  }, [fetchConfigs]);

  return (
    <ConfigContext.Provider
      value={{
        configs,
        loading,
        error,
        refreshConfigs: fetchConfigs,
        readConfig,
        saveConfig,
        deleteConfig,
      }}
    >
      {children}
    </ConfigContext.Provider>
  );
};

export const useConfig = () => {
  const context = useContext(ConfigContext);
  if (context === undefined) {
    throw new Error("useConfig must be used within a ConfigProvider");
  }
  return context;
};
