import type { Page } from "@playwright/test";

/**
 * Minimal structural shape of a db entry — only the fields the db app actually
 * reads (`DBCard`, `craftQuery`). Kept local rather than importing
 * `@gcsim/types` so the e2e harness stays dependency-free.
 */
export interface DbEntry {
	_id: string;
	create_date: number;
	config: string;
	description: string;
	submitter: string;
	accepted_tags: number[];
	summary: {
		mode: number;
		char_names: string[];
		target_count: number;
		mean_dps_per_target: number;
		team: { name: string }[];
	};
}

const team = (...names: string[]): { name: string }[] =>
	names.map((name) => ({ name }));

/**
 * The deterministic entries the stubbed `/api/db` returns for an unfiltered
 * browse. Only the first team contains Nahida, so a Nahida character filter
 * narrows the result set from two entries to one — an observable "list
 * updated" signal.
 */
export const dbEntries: DbEntry[] = [
	{
		_id: "aaaaaaaaaaaaaaaaaaaaaaaa",
		create_date: 1_700_000_000,
		config: "// nahida hyperbloom sample config\noptions iteration=1;\n",
		description: "Nahida hyperbloom sample",
		submitter: "e2e",
		accepted_tags: [8],
		summary: {
			mode: 0,
			char_names: ["nahida", "furina", "yelan", "raidenshogun"],
			target_count: 1,
			mean_dps_per_target: 123456,
			team: team("nahida", "furina", "yelan", "raidenshogun"),
		},
	},
	{
		_id: "bbbbbbbbbbbbbbbbbbbbbbbb",
		create_date: 1_700_000_100,
		config: "// hu tao vape sample config\noptions iteration=1;\n",
		description: "Hu Tao vape sample",
		submitter: "migrated",
		accepted_tags: [8],
		summary: {
			mode: 0,
			char_names: ["hutao", "yelan", "furina", "raidenshogun"],
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

// Stub GitHub "latest release" payload for the home page's <LatestVersion />,
// which fetches it directly (not via `/api`). Keeps the home test offline.
const GITHUB_RELEASE = { name: "v5.0", body: "e2e stub release notes" };

/**
 * The characters an included filter names in a `/api/db` query. `craftQuery`
 * serializes an included char as `"summary.char_names":"<char>"`; an excluded
 * char uses a `{ "$ne": ... }` object, which this pattern does not match.
 */
function includedCharsFromQuery(q: string): string[] {
	return [...q.matchAll(/"summary\.char_names":"(\w+)"/g)].map((m) => m[1]);
}

/**
 * Route every network dependency of the db app to a local, deterministic stub:
 *
 *  - `/api/db` returns {@link dbEntries}, filtered to the characters an included
 *    filter names (see {@link includedCharsFromQuery}) so a character filter
 *    visibly narrows the list;
 *  - `/api/assets/**` (avatars, weapons, misc art) returns a 1x1 PNG;
 *  - `api.github.com` (latest-release lookup) returns a fixed payload;
 *  - any other `/api/**` call returns an empty 200 so nothing reaches prod.
 *
 * The dev server proxies `/api` to production by default; these routes ensure
 * the spec never depends on live data and runs offline. Registered
 * catch-all-first so the specific handlers, added last, take precedence.
 */
export async function installDbRoutes(page: Page): Promise<void> {
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
	await page.route("**/api/db*", (route) => {
		const q = new URL(route.request().url()).searchParams.get("q") ?? "";
		const included = includedCharsFromQuery(q);
		const data =
			included.length > 0
				? dbEntries.filter((e) =>
						included.every((c) => e.summary.char_names.includes(c)),
					)
				: dbEntries;
		route.fulfill({
			contentType: "application/json",
			body: JSON.stringify({ data }),
		});
	});
}
