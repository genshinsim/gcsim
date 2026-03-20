import path from "node:path";
import { defineConfig, mergeConfig } from "vitest/config";
import baseConfig from "../../tooling/vitest/base.ts";

export default mergeConfig(
  baseConfig,
  defineConfig({
    resolve: {
      alias: {
        "@": path.resolve(import.meta.dirname, "src"),
      },
    },
    test: {
      setupFiles: ["./src/test-setup.ts"],
    },
  }),
);
