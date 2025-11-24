import React, { createContext, useContext, useState, useEffect } from 'react';
import { DiscoverConfigs } from '../../wailsjs/go/main/App';
import { config } from '../../wailsjs/go/config/models';

interface ConfigContextType {
  configs: config.ConfigItem[];
  loading: boolean;
  error: string | null;
  refreshConfigs: () => Promise<void>;
}

const ConfigContext = createContext<ConfigContextType | undefined>(undefined);

export const ConfigProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [configs, setConfigs] = useState<config.ConfigItem[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  const fetchConfigs = async () => {
    setLoading(true);
    setError(null);
    try {
      // DiscoverConfigs takes a projectPath. Empty string means discover in default locations.
      const items = await DiscoverConfigs("");
      setConfigs(items || []);
    } catch (err) {
      console.error("Failed to load configs:", err);
      setError(err instanceof Error ? err.message : "Failed to load configurations");
      // Mock data for development if Wails is not available (e.g. in browser)
      if (String(err).includes("window.go") || String(err).includes("is not a function")) {
         console.warn("Wails runtime not found. Using mock data.");
         // Optional: Set mock data here if we want to test UI without backend
      }
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchConfigs();
  }, []);

  return (
    <ConfigContext.Provider value={{ configs, loading, error, refreshConfigs: fetchConfigs }}>
      {children}
    </ConfigContext.Provider>
  );
};

export const useConfig = () => {
  const context = useContext(ConfigContext);
  if (context === undefined) {
    throw new Error('useConfig must be used within a ConfigProvider');
  }
  return context;
};
