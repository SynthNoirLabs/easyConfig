import type React from "react";
import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useState,
} from "react";
import { toast } from "sonner";
import {
  DeleteConfig,
  DiscoverConfigs,
  ReadConfig,
  SaveConfig,
} from "../../wailsjs/go/main/App";
import type { config } from "../../wailsjs/go/models";
import { EventsOn } from "../../wailsjs/runtime/runtime";
import {
  deleteDemoConfig,
  getDemoConfigs,
  isBrowserDemoMode,
  isWailsUnavailableError,
  readDemoConfig,
  saveDemoConfig,
} from "../mocks/browserDemoData";
import { debounce } from "../utils/debounce";

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
    if (isBrowserDemoMode()) {
      setConfigs(getDemoConfigs());
      setLoading(false);
      return;
    }
    try {
      // DiscoverConfigs takes a projectPath. Empty string means discover in default locations.
      const items = await DiscoverConfigs("");
      setConfigs(items || []);
    } catch (err) {
      if (isWailsUnavailableError(err)) {
        setConfigs(getDemoConfigs());
        setError(null);
        return;
      }
      console.error("Failed to load configs:", err);
      setError(
        err instanceof Error ? err.message : "Failed to load configurations",
      );
      toast.error(
        err instanceof Error
          ? err.message
          : "Failed to load configurations. Check backend logs.",
      );
    } finally {
      setLoading(false);
    }
  }, []);

  const readConfig = async (path: string): Promise<string> => {
    if (isBrowserDemoMode()) {
      return readDemoConfig(path);
    }
    try {
      return await ReadConfig(path);
    } catch (err) {
      if (isWailsUnavailableError(err)) {
        return readDemoConfig(path);
      }
      console.error("Failed to read config:", err);
      throw err;
    }
  };

  const saveConfig = async (path: string, content: string): Promise<void> => {
    if (isBrowserDemoMode()) {
      saveDemoConfig(path, content);
      return;
    }
    try {
      await SaveConfig(path, content);
    } catch (err) {
      if (isWailsUnavailableError(err)) {
        saveDemoConfig(path, content);
        return;
      }
      console.error("Failed to save config:", err);
      throw err;
    }
  };

  const deleteConfig = async (path: string): Promise<void> => {
    if (isBrowserDemoMode()) {
      deleteDemoConfig(path);
      await fetchConfigs();
      return;
    }
    try {
      await DeleteConfig(path);
      await fetchConfigs(); // Refresh list after delete
    } catch (err) {
      if (isWailsUnavailableError(err)) {
        deleteDemoConfig(path);
        await fetchConfigs();
        return;
      }
      console.error("Failed to delete config:", err);
      throw err;
    }
  };

  useEffect(() => {
    fetchConfigs();

    if (isBrowserDemoMode()) {
      return;
    }

    const debouncedRefresh = debounce(() => {
      toast.info("Configuration changed on disk. Refreshing…");
      void fetchConfigs();
    }, 300);

    const off = EventsOn("config:changed", debouncedRefresh);

    return () => {
      debouncedRefresh.cancel();
      off();
    };
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
