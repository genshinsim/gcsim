import { test as base } from "@playwright/test";
import { AppHarness } from "./app-harness";
import { DocsHarness } from "./docs-harness";

/**
 * Playwright test extended with two lazy per-test fixtures, each a fresh harness
 * bound to the test's page:
 *  - `app` — an {@link AppHarness} for the web app (`/simulator` + `/web`);
 *  - `docs` — a {@link DocsHarness} for the docs site.
 *
 * They are independent: a spec destructures only the one it needs, so a docs
 * spec never boots the wasm web server and vice versa. Specs use `test`/`expect`
 * from here instead of `@playwright/test` directly.
 */
export const test = base.extend<{ app: AppHarness; docs: DocsHarness }>({
	app: async ({ page }, use) => {
		await use(new AppHarness(page));
	},
	docs: async ({ page }, use) => {
		await use(new DocsHarness(page));
	},
});

export { expect } from "@playwright/test";
