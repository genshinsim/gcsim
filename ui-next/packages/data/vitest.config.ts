import { defineConfig, mergeConfig } from "vitest/config";
import baseConfig from "../../tooling/vitest/base.ts";

export default mergeConfig(baseConfig, defineConfig({}));
