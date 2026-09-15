import { test as base } from "@playwright/test";
import { AppHarness } from "./app-harness";
import { installDbRoutes } from "./db-fixtures";
import { DbHarness } from "./db-harness";
import { DocsHarness } from "./docs-harness";
import { installTaghelperRoutes } from "./taghelper-fixtures";
import { TaghelperHarness } from "./taghelper-harness";

/**
 * Playwright test extended with lazy per-test fixtures, each a fresh harness
 * bound to the test's page:
 *  - `app` — an {@link AppHarness} for the web app (`/simulator` + `/web`);
 *  - `docs` — a {@link DocsHarness} for the docs site;
 *  - `db` — a {@link DbHarness} for the db app, with its network stubbed
 *    (see {@link installDbRoutes}) before any navigation;
 *  - `taghelper` — a {@link TaghelperHarness} for the taghelper app, with its
 *    network stubbed (see {@link installTaghelperRoutes}) before any navigation.
 *
 * They are independent: a spec destructures only the one it needs, so a docs
 * spec never boots the wasm web server and vice versa. Specs use `test`/`expect`
 * from here instead of `@playwright/test` directly.
 */
export const test = base.extend<{
	app: AppHarness;
	docs: DocsHarness;
	db: DbHarness;
	taghelper: TaghelperHarness;
}>({
	app: async ({ page }, use) => {
		await use(new AppHarness(page));
	},
	docs: async ({ page }, use) => {
		await use(new DocsHarness(page));
	},
	db: async ({ page }, use) => {
		await installDbRoutes(page);
		await use(new DbHarness(page));
	},
	taghelper: async ({ page }, use) => {
		await installTaghelperRoutes(page);
		await use(new TaghelperHarness(page));
	},
});

export { expect } from "@playwright/test";
