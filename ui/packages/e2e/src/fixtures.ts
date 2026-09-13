import { test as base } from "@playwright/test";
import { AppHarness } from "./app-harness";

/**
 * Playwright test extended with an `app` fixture — a fresh {@link AppHarness}
 * bound to the test's page. Specs use `test`/`expect` from here instead of
 * `@playwright/test` directly.
 */
export const test = base.extend<{ app: AppHarness }>({
	app: async ({ page }, use) => {
		await use(new AppHarness(page));
	},
});

export { expect } from "@playwright/test";
