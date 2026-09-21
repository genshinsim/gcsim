import { expect, type Locator, type Page } from "@playwright/test";

export class DbHomePage {
	readonly page: Page;
	readonly browse: Locator;
	readonly welcome: Locator;
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

	async waitForLoaded(): Promise<void> {
		await expect(this.welcome).toBeVisible();
		await expect(this.browse).toBeVisible();
		await expect(this.whatsNew).toBeVisible();
	}
}
