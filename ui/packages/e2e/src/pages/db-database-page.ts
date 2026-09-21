import { expect, type Locator, type Page } from "@playwright/test";

/**
 * Drives the db app's browse route (`/database`): the sticky command bar (custom
 * filter search, filter sheet trigger, sort control, live result count), the
 * full entry cards, and the filter sheet.
 *
 * Locators ride observable contracts — accessible names, visible copy, and ARIA
 * roles — since the app ships no data-testids.
 */
export class DbDatabasePage {
	readonly page: Page;
	/** The command-bar custom-filter search input (placeholder "Custom Filter"). */
	readonly searchBox: Locator;
	/** The command-bar button that opens the filter sheet (accessible name "Filter"). */
	readonly filterButton: Locator;
	/** "Copy config" button in each entry card — one per entry. */
	readonly copyConfigButtons: Locator;
	/** The "Open in viewer" link in each entry card — links cross-app. */
	readonly openInViewerLinks: Locator;
	/** Filter sheet collapsible section headers. */
	readonly charactersSection: Locator;
	readonly tagsSection: Locator;
	readonly sortBySection: Locator;
	/** The character search input inside the expanded Characters section. */
	readonly charSearch: Locator;
	/** Character portrait images in the expanded Characters section / cards. */
	readonly characterPortraits: Locator;

	constructor(page: Page) {
		this.page = page;
		this.searchBox = page.getByPlaceholder("Custom Filter");
		this.filterButton = page.getByRole("button", {
			name: "Filter",
			exact: true,
		});
		this.copyConfigButtons = page.getByRole("button", { name: "Copy config" });
		this.openInViewerLinks = page.getByRole("link", {
			name: "Open in viewer",
		});
		this.charactersSection = page.getByRole("button", { name: /Characters/ });
		this.tagsSection = page.getByRole("button", { name: /Tags/ });
		this.sortBySection = page.getByRole("button", { name: /Sort by/ });
		this.charSearch = page.getByPlaceholder("Type to search...");
		this.characterPortraits = page.locator('img[src^="/api/assets/avatar/"]');
	}

	/** Navigate to the browse route and wait for React to mount into `#root`. */
	async goto(): Promise<void> {
		await this.page.goto("/database");
		await expect(this.page.locator("#root")).not.toBeEmpty();
	}

	/**
	 * Assert the browse view rendered: the result count, the search box, the
	 * filter trigger, and at least one entry card (its Copy config / Open in
	 * viewer controls). Structural only.
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

	/** Open the filter sheet and wait for its Characters section to appear. */
	async openFilterPanel(): Promise<void> {
		await this.filterButton.click();
		await expect(this.charactersSection).toBeVisible();
	}

	/** Expand the Characters section and wait for its portrait picker to render. */
	async expandCharacters(): Promise<void> {
		await this.charactersSection.click();
		await expect(this.characterPortraits.first()).toBeVisible();
	}

	/**
	 * Filter to a single character: open the sheet, expand Characters, type the
	 * name into the picker's search, pick the match, then close the sheet.
	 * Dispatches an include filter, which refetches `/api/db` with the narrowed
	 * query.
	 */
	async filterByCharacter(name: string): Promise<void> {
		await this.openFilterPanel();
		await this.expandCharacters();
		await this.charSearch.fill(name);
		await this.page.getByRole("button", { name, exact: true }).first().click();
		await this.page.keyboard.press("Escape");
	}
}
