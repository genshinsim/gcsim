import type { Page } from "@playwright/test";

/**
 * Deterministic fixtures for the taghelper smoke spec.
 *
 * The taghelper dev server proxies `/api` to production (simimpact.app), so the
 * spec must not read live data. {@link installTaghelperRoutes} intercepts every
 * `/api` request at the browser and answers from these fixtures, making the run
 * hermetic and repeatable.
 */

/**
 * Minimal structural shape of a taghelper db entry — only the fields the app
 * actually reads (`DBCard`, the "same team" query). Kept local rather than
 * importing `@gcsim/types` so the e2e harness stays dependency-free.
 */
export interface TaghelperEntry {
	_id: string;
	share_key: string;
	submitter: string;
	description: string;
	create_date: number;
	accepted_tags: number[];
	summary: {
		mode: number;
		target_count: number;
		mean_dps_per_target: number;
		sim_duration: { min: number; max: number; mean: number; sd: number };
		char_names: string[];
		team: { name: string; element: string; level: number; cons: number }[];
	};
}

/** Route param the spec navigates to (`/id/:id`). */
export const MAIN_ID = "smoke-main";

// The entry's source tag. DBCard renders `tags.json[SOURCE_TAG_ID].display_name`
// as a chip; keeping the id and its rendered name together here is the single
// source of truth the spec asserts against.
const SOURCE_TAG_ID = 1;
export const SOURCE_TAG_NAME = "gcsim";

const CHARS = ["raiden", "nahida", "kazuha", "bennett"] as const;

/**
 * The entry rendered by the main card at `/id/:id`. `summary.char_names` is
 * non-empty so the app fires the "same characters" query the spec also stubs.
 */
export const mainEntry: TaghelperEntry = {
	_id: MAIN_ID,
	share_key: "smokeshare",
	submitter: "smoke-tester",
	description: "smoke fixture entry",
	create_date: 1_700_000_000,
	accepted_tags: [SOURCE_TAG_ID],
	summary: {
		mode: 1,
		target_count: 1,
		mean_dps_per_target: 123_456,
		sim_duration: { min: 90, max: 90, mean: 90, sd: 0 },
		char_names: [...CHARS],
		team: CHARS.map((name, i) => ({
			name,
			element: "electro",
			level: 90,
			cons: i,
		})),
	},
};

/** One related entry with a distinct `_id`, so an "existing sims" row renders. */
export const relatedEntry: TaghelperEntry = {
	...mainEntry,
	_id: "smoke-related",
	share_key: "relatedshare",
	description: "related fixture entry",
};

// 1×1 transparent PNG — every `/api/assets/...` request resolves to this so the
// suite never fetches portrait/weapon images from production.
const PNG_1x1 = Buffer.from(
	"iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNk+M8AAAMCAoLl9wAAAABJRU5ErkJggg==",
	"base64",
);

/** Intercept all `/api` traffic and answer from the fixtures above. */
export async function installTaghelperRoutes(
	page: Page,
	opts: { related?: TaghelperEntry[] } = {},
): Promise<void> {
	const related = opts.related ?? [relatedEntry];
	await page.route(
		(url) => url.pathname.startsWith("/api/db/id/"),
		(route) => route.fulfill({ json: mainEntry }),
	);
	await page.route(
		(url) => url.pathname === "/api/db",
		(route) => route.fulfill({ json: { data: related } }),
	);
	await page.route(
		(url) => url.pathname.startsWith("/api/assets/"),
		(route) => route.fulfill({ contentType: "image/png", body: PNG_1x1 }),
	);
}
