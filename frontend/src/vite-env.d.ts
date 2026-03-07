/// <reference types="vite/client" />

interface Window {
  go?: {
    main: {
      App: Record<string, unknown>;
    };
  };
  runtime?: Record<string, unknown>;
}
