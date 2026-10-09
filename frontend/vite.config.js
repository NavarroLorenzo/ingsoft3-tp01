import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import { configDefaults } from "vitest/config";

export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      "/api": "http://localhost:8080",
      "/health": "http://localhost:8080"
    }
  },
  test: {
    environment: "node",
    include: ["tests/**/*.test.js"],
    exclude: [...configDefaults.exclude, "e2e/**"],
    coverage: {
      provider: "v8",
      reporter: ["text", "json-summary", "html", "lcov"],
      // Docker monta coverage/; Vitest puede limpiar esta subcarpeta sin borrar el montaje.
      reportsDirectory: "coverage/report",
      include: [
        "src/api/**/*.js",
        "src/services/**/*.js",
        "src/utils/**/*.js"
      ],
      thresholds: {
        lines: 80,
        functions: 80,
        statements: 80,
        branches: 80
      }
    }
  }
});
