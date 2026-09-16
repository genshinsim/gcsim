import { expect, type Locator, type Page } from "@playwright/test";

/**
 * Drives the dash home route (`/`): the top nav bar, the "Get started" CTA, and
 * the featured-submission card.
 *
 * Locators ride observable contracts (accessible names, hardcoded button text)
 * since the app ships no data-testids. The featured card only renders once
 * `/api/db` resolves, so a spec must stub that route (see `installDbRoutes`)
 * before navigating.
 */
export class DashPage {
	readonly page: Page;
	/** "Get started" CTA linking to `/simulator`. */
	readonly getStarted: Locator;
	/** Nav-bar entry to the simulator. */
	readonly navSimulator: Locator;
	/** "Show Detail" button rendered on each featured submission card. */
	readonly showDetail: Locator;
	/** "Visit the Teams DB" CTA below the featured card (a Blueprint
	 * AnchorButton, which carries role="button"). */
	readonly visitTeamsDb: Locator;

	constructor(page: Page) {
		this.page = page;
		this.getStarted = page.getByRole("button", { name: "Get started" });
		// The nav renders a desktop and a (collapsed) mobile copy; take the first.
		this.navSimulator = page.getByRole("button", { name: "Simulator" }).first();
		this.showDetail = page.getByRole("link", { name: "Show Detail" });
		this.visitTeamsDb = page.getByRole("button", {
			name: "Visit the Teams DB",
		});
	}

	/** Navigate to the home route and wait for React to mount into `#root`. */
	async goto(): Promise<void> {
		await this.page.goto("/");
		await expect(this.page.locator("#root")).not.toBeEmpty();
	}

	/**
	 * Assert the dash chrome rendered: title, the nav bar's Simulator entry, and
	 * the "Get started" CTA. Structural only.
	 */
	async waitForLoaded(): Promise<void> {
		await expect(this.page).toHaveTitle("gcsim - simulation impact");
		await expect(this.navSimulator).toBeVisible();
		await expect(this.getStarted).toBeVisible();
	}

	/**
	 * Assert the featured-submission section rendered: at least one card with its
	 * "Show Detail" link, and the "Visit the Teams DB" CTA. Depends on a stubbed
	 * `/api/db`.
	 */
	async waitForFeatured(): Promise<void> {
		await expect(this.showDetail.first()).toBeVisible();
		await expect(this.visitTeamsDb).toBeVisible();
	}
}
