import tailwindcss from "@tailwindcss/vite";
import react from "@vitejs/plugin-react";
import * as path from "path";
import { defineConfig } from "vite";

// https://vitejs.dev/config/
export default defineConfig({
	plugins: [tailwindcss(), react()],
	server: {
		proxy: {
			"/api": {
				target: "https://gcsim.app",
				changeOrigin: true,
			},
		},
	},
	resolve: {
		alias: [{ find: "@", replacement: path.resolve(import.meta.dirname, "src") }],
	},
});
