import { expect, type Locator, type Page } from "@playwright/test";

export class DashPage {
	readonly page: Page;
	readonly openSimulator: Locator;
	readonly navSimulator: Locator;
	readonly sharedByOthersHeading: Locator;
	readonly teamCards: Locator;

	constructor(page: Page) {
		this.page = page;
		this.openSimulator = page.getByRole("link", { name: "Open Simulator" });
		// The nav renders a desktop and a (collapsed) mobile copy; take the first.
		this.navSimulator = page
			.getByRole("link", { name: "Simulator", exact: true })
			.first();
		this.sharedByOthersHeading = page.getByRole("heading", {
			level: 2,
			name: "Shared by others, courtesy of KQM DB",
		});
		this.teamCards = page.locator("a[href^='https://db.kqm.gg/db/']");
	}

	/** Navigate to the home route and wait for React to mount into `#root`. */
	async goto(): Promise<void> {
		await this.page.goto("/");
		await expect(this.page.locator("#root")).not.toBeEmpty();
	}

	async waitForLoaded(): Promise<void> {
		await expect(this.page).toHaveTitle("gcsim - simulation impact");
		await expect(this.navSimulator).toBeVisible();
		await expect(this.openSimulator).toBeVisible();
	}

	/** Depends on a stubbed KQM DB `/api/db` (see `installOfflineRoutes`). */
	async waitForSharedByOthers(): Promise<void> {
		await expect(this.sharedByOthersHeading).toBeVisible();
		await expect(this.teamCards.first()).toBeVisible();
	}
}
