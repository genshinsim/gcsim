import react from "@vitejs/plugin-react";
import * as git from "git-rev-sync";
import { visualizer } from "rollup-plugin-visualizer";
import { type ConfigEnv, defineConfig } from "vite";
import tsconfigPaths from "vite-tsconfig-paths";

export default (_env: ConfigEnv) => {
	process.env.VITE_GIT_COMMIT_HASH = git.long();
	process.env.VITE_GIT_BRANCH = git.branch();

	return defineConfig({
		plugins: [react(), tsconfigPaths(), visualizer()],
		build: {
			rollupOptions: {
				output: {
					// TODO: remove after using react lazy imports?
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
					},
				},
			},
		},
		server: {
			proxy: {
				"/api": {
					target: "https://gcsim.app",
					changeOrigin: true,
				},
				"/hastebin/post": {
					target: "https://hastebin.com/documents",
					changeOrigin: true,
					rewrite: (path) => path.replace(/^\/hastebin\/post/, ""),
				},
				"/hastebin/get": {
					target: "https://hastebin.com/raw/",
					changeOrigin: true,
					rewrite: (path) => path.replace(/^\/hastebin\/get/, ""),
				},
			},
		},
	});
};
