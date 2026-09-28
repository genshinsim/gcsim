import { expect, type Locator, type Page } from "@playwright/test";

/**
 * Drives the dash home route (`/`): the top nav bar, the "Open Simulator" CTA,
 * and the "Shared by others" section.
 *
 * Locators ride observable contracts (ARIA roles and accessible names) since
 * the app ships no data-testids. The shared-by-others cards only render once
 * `/api/db` resolves, so a spec must stub that route (see `installDbRoutes`)
 * before navigating.
 */
export class DashPage {
	readonly page: Page;
	/** "Open Simulator" hero CTA linking to `/simulator`. */
	readonly openSimulator: Locator;
	/** Nav-bar link to the simulator. */
	readonly navSimulator: Locator;
	/** "Shared by others" section heading. */
	readonly sharedHeading: Locator;
	/** Team cards in the "Shared by others" section. Each card is a link whose
	 * accessible name ends with its fixed "DPS / target" label. */
	readonly teamCards: Locator;

	constructor(page: Page) {
		this.page = page;
		this.openSimulator = page.getByRole("link", { name: "Open Simulator" });
		// The nav renders a desktop and a (collapsed) mobile copy; take the first.
		this.navSimulator = page
			.getByRole("link", { name: "Simulator", exact: true })
			.first();
		this.sharedHeading = page.getByRole("heading", {
			level: 2,
			name: "Shared by others",
		});
		this.teamCards = page.getByRole("link", { name: /DPS \/ target/ });
	}

	/** Navigate to the home route and wait for React to mount into `#root`. */
	async goto(): Promise<void> {
		await this.page.goto("/");
		await expect(this.page.locator("#root")).not.toBeEmpty();
	}

	/**
	 * Assert the dash chrome rendered: title, the nav bar's Simulator link, and
	 * the "Open Simulator" CTA. Structural only.
	 */
	async waitForLoaded(): Promise<void> {
		await expect(this.page).toHaveTitle("gcsim - simulation impact");
		await expect(this.navSimulator).toBeVisible();
		await expect(this.openSimulator).toBeVisible();
	}

	/**
	 * Assert the "Shared by others" section rendered: its heading and at least
	 * one team card. Depends on a stubbed `/api/db`.
	 */
	async waitForSharedByOthers(): Promise<void> {
		await expect(this.sharedHeading).toBeVisible();
		await expect(this.teamCards.first()).toBeVisible();
	}
}
