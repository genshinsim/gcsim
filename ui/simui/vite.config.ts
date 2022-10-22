import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import { resolve } from "path";
import dts from "vite-plugin-dts";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    react(),
    dts({
      include: "lib",
      exclude: "src",
    }),
  ],
  resolve: {
    alias: {
      "~": resolve(__dirname, "lib"),
    },
  },
  build: {
    lib: {
      entry: resolve(__dirname, "lib/index.ts"),
      formats: ["es"],
    },
    rollupOptions: {
      external: ["react", "react-dom", "axios"],
      output: {
        globals: {
          react: "React",
          "react-dom": "ReactDOM",
          "axios": "axios",
        },
      },
    },
  },
});
