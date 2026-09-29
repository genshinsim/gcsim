import type { Page } from "@playwright/test";

const team = (...names: string[]): { name: string }[] =>
	names.map((name) => ({ name }));

/** The deterministic entries the stubbed KQM DB `/api/db` returns. */
export const kqmEntries = [
	{
		_id: "aaaaaaaaaaaa",
		create_date: "1700000000",
		description: "Nahida hyperbloom sample",
		submitter: "e2e",
		summary: {
			mode: 2,
			target_count: 1,
			mean_dps_per_target: 123456,
			team: team("nahida", "furina", "yelan", "raidenshogun"),
		},
	},
	{
		_id: "bbbbbbbbbbbb",
		create_date: "1700000100",
		description: "Hu Tao vape sample",
		submitter: "e2e",
		summary: {
			mode: 1,
			target_count: 1,
			mean_dps_per_target: 234567,
			team: team("hutao", "yelan", "furina", "raidenshogun"),
		},
	},
];

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
 *  - `db.kqm.gg/api/db` (Shared by others) returns {@link kqmEntries};
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
	await page.route("https://db.kqm.gg/api/db*", (route) =>
		route.fulfill({
			contentType: "application/json",
			body: JSON.stringify({ data: kqmEntries }),
		}),
	);
}
