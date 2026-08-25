import { defineConfig } from '@playwright/test';

export default defineConfig({
  testDir: './tests',
  timeout: 30000,
  retries: 1,
  use: {
    headless: true,
    launchOptions: {
      executablePath: process.env.PLAYWRIGHT_BROWSER_PATH || undefined,
      args: ['--no-sandbox', '--single-process', '--disable-gpu', '--no-zygote', '--disable-dev-shm-usage'],
    },
    baseURL: 'http://localhost:4173',
  },
  webServer: {
    command: 'npm run preview',
    port: 4173,
    reuseExistingServer: !process.env.CI,
    timeout: 30000,
  },
});
