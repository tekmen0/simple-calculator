import react from "@vitejs/plugin-react";
import { defineConfig } from "vitest/config";

export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      "/api/": "http://127.0.0.1:8080",
    },
  },
  test: {
    coverage: {
      include: ["src/**/*.{ts,tsx}"],
      exclude: ["src/main.tsx", "src/**/*.test.{ts,tsx}"],
    },
  },
});
