import { defineConfig, devices } from "@playwright/test";

const PORT = 5273;
const HOST = `http://localhost:${PORT}`;

/**
 * Chromium-only e2e config for the db ("Simpact") app. Unlike the web app, the
 * db app is pure front-end — no wasm — so this boots only the `vite` dev server
 * (fast, no Go toolchain needed). The spec stubs every `/api` call, so the dev
 * server's default proxy to production is never exercised.
 */
export default defineConfig({
	testDir: "./tests/db",
	fullyParallel: true,
	forbidOnly: !!process.env.CI,
	retries: process.env.CI ? 1 : 0,
	workers: 1,
	timeout: 60_000,
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
		// `exec vite` (not the `dev` script) so `--port` reaches vite directly:
		// the `dev` script's port default is 5173, and args after pnpm's `--` are
		// swallowed, which would leave the server off 5273.
		command:
			"pnpm --filter @gcsim/db exec vite --host --port 5273 --strictPort",
		url: HOST,
		reuseExistingServer: !process.env.CI,
		// The dev server's first compile can be slow on a cold cache.
		timeout: 120_000,
		stdout: "pipe",
		stderr: "pipe",
	},
});
