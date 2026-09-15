import { expect, type Locator, type Page } from "@playwright/test";

/**
 * Drives the db app's home route (`/`): the welcome copy and the "Get started"
 * call to action that routes into `/database`.
 *
 * The app ships no data-testids, so locators ride accessible names and visible
 * copy (the English localization strings).
 */
export class DbHomePage {
	readonly page: Page;
	/** Blueprint button, accessible name "Get started" — routes to /database. */
	readonly getStarted: Locator;
	/** The welcome card heading ("Welcome to Simpact"). */
	readonly welcome: Locator;
	/** The tag-list header that introduces the tag descriptions. */
	readonly tagsCopy: Locator;

	constructor(page: Page) {
		this.page = page;
		this.getStarted = page.getByRole("button", { name: "Get started" });
		this.welcome = page.getByText("Welcome to Simpact");
		this.tagsCopy = page.getByText("Below are the current available tags:");
	}

	/** Navigate to the home page and wait for React to mount into `#root`. */
	async goto(): Promise<void> {
		await this.page.goto("/");
		await expect(this.page.locator("#root")).not.toBeEmpty();
	}

	/**
	 * Assert the home page rendered: welcome copy, the tag-list copy, and the
	 * "Get started" CTA.
	 */
	async waitForLoaded(): Promise<void> {
		await expect(this.welcome).toBeVisible();
		await expect(this.tagsCopy).toBeVisible();
		await expect(this.getStarted).toBeVisible();
	}
}
