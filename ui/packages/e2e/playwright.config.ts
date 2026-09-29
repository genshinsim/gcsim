import { defineConfig, devices } from "@playwright/test";

const PORT = 5173;
const HOST = `http://localhost:${PORT}`;

/**
 * Chromium e2e config (plus Firefox and WebKit for the wasm CSP spec). Runs
 * against the *dev server* (build wasm, then
 * `vite`), because a production build points the wasm URL at a remote origin
 * (R2 / `/api/wasm/...`) that is not available locally.
 *
 * The `webServer` builds the wasm binary first (needs Go + `task` on PATH — see
 * README) and then serves the app on a fixed, strict port so the harness knows
 * where to connect.
 */
export default defineConfig({
	testDir: "./tests",
	testIgnore: ["**/docs/**"],
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
		// Browsers differ in how the page's CSP applies to wasm, which Chromium
		// alone doesn't catch; only the CSP spec runs in these.
		{
			name: "firefox",
			use: { ...devices["Desktop Firefox"] },
			testMatch: /wasm-csp\.spec\.ts/,
		},
		{
			name: "webkit",
			use: { ...devices["Desktop Safari"] },
			testMatch: /wasm-csp\.spec\.ts/,
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
