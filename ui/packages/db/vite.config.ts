import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import tsconfigPaths from 'vite-tsconfig-paths';
import { visualizer } from "rollup-plugin-visualizer";

export default defineConfig({
  plugins: [
    react(),
    tsconfigPaths(),
    visualizer()
  ],
  build: {
    rollupOptions: {
      output: {
        manualChunks: (id) => {
          if (id.includes("node_modules")) {
            if (id.includes("@blueprintjs") && id.includes("icons")) {
              return "blueprint-icons";
            }
            if (id.includes("prismjs") || id.includes("pako")) {
              return "core";
            }
            return "vendor";
          }
          return "core";
        }
      }
    }
  },
  server: {
    proxy: {
      "/api": {
        target: "https://gcsim.app",
        changeOrigin: true
      },
    }
  }
});