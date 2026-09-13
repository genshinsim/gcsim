import { defineConfig, devices } from "@playwright/test";

const PORT = 5173;
const HOST = `http://localhost:${PORT}`;

/**
 * Chromium-only e2e config. Runs against the *dev server* (build wasm, then
 * `vite`), because a production build points the wasm URL at a remote origin
 * (R2 / `/api/wasm/...`) that is not available locally.
 *
 * The `webServer` builds the wasm binary first (needs Go + `task` on PATH — see
 * README) and then serves the app on a fixed, strict port so the harness knows
 * where to connect.
 */
export default defineConfig({
	testDir: "./tests",
	fullyParallel: true,
	forbidOnly: !!process.env.CI,
	retries: process.env.CI ? 1 : 0,
	workers: 1,
	// wasm boot + a sim run comfortably exceed Playwright's 30s default.
	timeout: 120_000,
	expect: { timeout: 15_000 },
	reporter: [["html", { open: "never" }], ["list"]],
	use: {
		baseURL: HOST,
		trace: "retain-on-failure",
		screenshot: "only-on-failure",
		video: "retain-on-failure",
	},
	projects: [
		{
			name: "chromium",
			use: { ...devices["Desktop Chrome"] },
		},
	],
	webServer: {
		command:
			"pnpm --filter @gcsim/executors build:wasm:web && pnpm --filter @gcsim/web dev -- --port 5173 --strictPort",
		url: HOST,
		reuseExistingServer: !process.env.CI,
		// wasm build (go build) plus the dev server's first compile can be slow.
		timeout: 240_000,
		stdout: "pipe",
		stderr: "pipe",
	},
});
