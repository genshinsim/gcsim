import tailwindcss from "@tailwindcss/vite";
import react from "@vitejs/plugin-react-swc";
import { visualizer } from "rollup-plugin-visualizer";
import { defineConfig } from "vite";

export default defineConfig({
	plugins: [tailwindcss(), react(), visualizer()],
	resolve: { tsconfigPaths: true },
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
		},
	},
});
