import { defineConfig } from "@playwright/test";
import { FRONT } from "./e2e/entorno.js";

export default defineConfig({
  testDir: "./e2e",
  timeout: 60_000,
  expect: { timeout: 15_000 },
  workers: 1,
  forbidOnly: Boolean(process.env.CI),
  retries: 1,
  reporter: [["html", { open: "never" }], ["list"]],
  use: {
    baseURL: FRONT,
    browserName: "chromium",
    trace: "on-first-retry",
    screenshot: "only-on-failure"
  }
});
