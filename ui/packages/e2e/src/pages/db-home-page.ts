import { expect, type Locator, type Page } from "@playwright/test";

/**
 * Drives the db app's home route (`/`): the Gauge hero with its "Browse
 * database" CTA (routes into `/database`) and the shared "What's new" feed.
 *
 * The app ships no data-testids, so locators ride accessible names and visible
 * copy (the English localization strings).
 */
export class DbHomePage {
	readonly page: Page;
	/** Primary CTA, accessible name "Browse database" — routes to /database. */
	readonly browse: Locator;
	/** The hero heading ("Welcome to Simpact"). */
	readonly welcome: Locator;
	/** The "What's new" section heading over the shared release feed. */
	readonly whatsNew: Locator;

	constructor(page: Page) {
		this.page = page;
		this.browse = page.getByRole("button", { name: "Browse database" });
		this.welcome = page.getByRole("heading", { name: "Welcome to Simpact" });
		this.whatsNew = page.getByRole("heading", { name: "What's new" });
	}

	/** Navigate to the home page and wait for React to mount into `#root`. */
	async goto(): Promise<void> {
		await this.page.goto("/");
		await expect(this.page.locator("#root")).not.toBeEmpty();
	}

	/**
	 * Assert the home page rendered: the hero heading, the "Browse database" CTA,
	 * and the "What's new" section.
	 */
	async waitForLoaded(): Promise<void> {
		await expect(this.welcome).toBeVisible();
		await expect(this.browse).toBeVisible();
		await expect(this.whatsNew).toBeVisible();
	}
}
