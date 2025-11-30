import { defineConfig } from '@playwright/test';

export default defineConfig({
  timeout: 60_000,
  testDir: './tests',
  fullyParallel: true,
  reporter: [['list']],
  use: {
    baseURL: 'http://localhost:4173',
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
    video: 'retain-on-failure',
  },
  webServer: {
    command: 'npm run preview -- --host --port 4173',
    url: 'http://localhost:4173',
    timeout: 60_000,
    reuseExistingServer: !process.env.CI,
  },
});
