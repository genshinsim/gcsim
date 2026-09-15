import { defineConfig, devices } from "@playwright/test";

const PORT = 4173;
const HOST = `http://localhost:${PORT}`;

/**
 * Chromium-only e2e config for the **docs** site (`@gcsim/docs`, Docusaurus).
 *
 * Unlike the web config, this serves a *production build* (`docusaurus build`
 * then `docusaurus serve`): docs has no wasm/workers/backend, and building first
 * makes the smoke spec catch build-time breakage too — broken links (config's
 * `onBrokenLinks: "throw"`), MDX compile errors, a bad sidebar entry.
 */
export default defineConfig({
	testDir: "./tests/docs",
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
		command:
			"pnpm --filter @gcsim/docs build && pnpm --filter @gcsim/docs serve --port 4173 --no-open",
		url: HOST,
		reuseExistingServer: !process.env.CI,
		// The Docusaurus production build can be slow on a cold cache.
		timeout: 240_000,
		stdout: "pipe",
		stderr: "pipe",
	},
});
