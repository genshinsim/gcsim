import { cloudflareTest } from "@cloudflare/vitest-pool-workers";
import { defineConfig } from "vitest/config";

const testKey = "000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f";

export default defineConfig({
	plugins: [
		cloudflareTest({
			wrangler: { configPath: "./wrangler.jsonc" },
			miniflare: {
				compatibilityDate: "2026-08-22",
				assets: {
					directory: "./test/site",
					binding: "ASSETS",
					routerConfig: { invoke_user_worker_ahead_of_assets: true },
				},
				bindings: {
					API_ENDPOINT: "https://backend.test",
					SHARE_KEYS: JSON.stringify({
						k3: { class: "prod", key: testKey },
						k4: { class: "dev", key: testKey },
					}),
				},
			},
		}),
	],
});
