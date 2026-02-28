const { defineConfig } = require("@playwright/test");

module.exports = defineConfig({
  testDir: "./e2e",
  timeout: 30000,
  use: {
    baseURL: "http://127.0.0.1:3001"
  },
  webServer: {
    command: "npm run dev",
    port: 3001,
    reuseExistingServer: true
  }
});
