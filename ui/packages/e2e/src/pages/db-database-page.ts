import { expect, type Locator, type Page } from "@playwright/test";

/**
 * Drives the db app's browse route (`/database`): the action bar (search box,
 * filter funnel, result count), the entry cards, and the filter drawer.
 *
 * Locators ride observable contracts — accessible names, visible copy, and
 * Blueprint classes — since the app ships no data-testids.
 */
export class DbDatabasePage {
	readonly page: Page;
	/** The character MultiSelect search box (placeholder "Type to search..."). */
	readonly searchBox: Locator;
	/** The funnel button that opens the filter drawer (48x48 primary button). */
	readonly filterButton: Locator;
	/** "Copy Config" button in each entry card's footer — one per entry. */
	readonly copyConfigButtons: Locator;
	/** The `<a>` wrapping each "Open in Viewer" button — links cross-app. */
	readonly openInViewerLinks: Locator;
	/** Filter drawer section headers. */
	readonly charactersSection: Locator;
	readonly tagsSection: Locator;
	readonly sortBySection: Locator;

	constructor(page: Page) {
		this.page = page;
		this.searchBox = page.getByPlaceholder("Type to search...");
		this.filterButton = page.locator(
			"button.bp4-button.bp4-intent-primary.w-12.h-12",
		);
		this.copyConfigButtons = page.getByRole("button", { name: "Copy Config" });
		this.openInViewerLinks = page.locator("a", {
			has: page.getByRole("button", { name: "Open in Viewer" }),
		});
		this.charactersSection = page.getByRole("button", { name: /Characters/ });
		this.tagsSection = page.getByRole("button", { name: /Tags/ });
		this.sortBySection = page.getByRole("button", { name: /Sort by/ });
	}

	/** Navigate to the browse route and wait for React to mount into `#root`. */
	async goto(): Promise<void> {
		await this.page.goto("/database");
		await expect(this.page.locator("#root")).not.toBeEmpty();
	}

	/**
	 * Assert the browse view rendered: the result count, the search box, the
	 * filter funnel, and at least one entry card (its Copy Config / Open in
	 * Viewer controls). Structural only.
	 */
	async waitForBrowse(): Promise<void> {
		await this.expectShowing(2);
		await expect(this.searchBox).toBeVisible();
		await expect(this.filterButton).toBeVisible();
		await expect(this.copyConfigButtons.first()).toBeVisible();
		await expect(this.openInViewerLinks.first()).toBeVisible();
	}

	/** Assert the "Showing N simulations" count reads exactly `n`. */
	async expectShowing(n: number): Promise<void> {
		await expect(this.page.getByText(`Showing ${n} simulations`)).toBeVisible();
	}

	/** Open the filter drawer and wait for its Characters section to appear. */
	async openFilterPanel(): Promise<void> {
		await this.filterButton.click();
		await expect(this.charactersSection).toBeVisible();
	}

	/**
	 * Filter to a single character via the search box: type its name, then pick
	 * it from the suggestion menu. Dispatches an include filter, which refetches
	 * `/api/db` with the narrowed query.
	 */
	async filterByCharacter(name: string): Promise<void> {
		await this.searchBox.click();
		await this.searchBox.pressSequentially(name);
		await this.page
			.locator(".bp4-menu-item", { hasText: name })
			.first()
			.click();
	}
}
