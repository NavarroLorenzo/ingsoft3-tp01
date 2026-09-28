import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

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
        lines: 100,
        functions: 100,
        statements: 100,
        branches: 100
      }
    }
  }
});
