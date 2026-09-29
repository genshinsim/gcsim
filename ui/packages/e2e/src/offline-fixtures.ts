import type { Page } from "@playwright/test";

// A 1x1 transparent PNG. Every avatar/asset request is fulfilled with this so
// the spec renders offline and never reaches production `/api/assets`.
const PNG_1x1 = Buffer.from(
	"iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==",
	"base64",
);

// Stub GitHub "latest release" payload, fetched directly (not via `/api`).
const GITHUB_RELEASE = { name: "v5.0", body: "e2e stub release notes" };

/**
 * Route the web app's network to a local, deterministic stub:
 *
 *  - `/api/assets/**` (avatars, weapons, misc art) returns a 1x1 PNG;
 *  - `api.github.com` (latest-release lookup) returns a fixed payload;
 *  - any other `/api/**` call returns an empty 200 so nothing reaches prod.
 *
 * The dev server proxies `/api` to production by default; these routes ensure
 * the spec never depends on live data and runs offline. Registered
 * catch-all-first so the specific handlers, added last, take precedence.
 */
export async function installOfflineRoutes(page: Page): Promise<void> {
	await page.route("**/api/**", (route) =>
		route.fulfill({ status: 200, contentType: "application/json", body: "{}" }),
	);
	await page.route("**/api/assets/**", (route) =>
		route.fulfill({ contentType: "image/png", body: PNG_1x1 }),
	);
	await page.route("https://api.github.com/**", (route) =>
		route.fulfill({
			contentType: "application/json",
			body: JSON.stringify(GITHUB_RELEASE),
		}),
	);
}
