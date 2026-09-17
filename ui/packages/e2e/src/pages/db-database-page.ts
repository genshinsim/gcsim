import { expect, type Locator, type Page } from "@playwright/test";

/**
 * Drives the db app's browse route (`/database`): the action bar (search box,
 * filter funnel, result count), the entry cards, and the filter drawer.
 *
 * Locators ride observable contracts — accessible names, visible copy, and ARIA
 * roles — since the app ships no data-testids.
 */
export class DbDatabasePage {
	readonly page: Page;
	/**
	 * The character MultiSelect trigger, labelled "Type to search...". Clicking
	 * it opens the cmdk popover holding the actual search input.
	 */
	readonly searchBox: Locator;
	/** The funnel button that opens the filter drawer (aria-label "Filter"). */
	readonly filterButton: Locator;
	/** "Copy Config" button in each entry card's footer — one per entry. */
	readonly copyConfigButtons: Locator;
	/** The `<a>` wrapping each "Open in Viewer" button — links cross-app. */
	readonly openInViewerLinks: Locator;
	/** Filter drawer section headers. */
	readonly charactersSection: Locator;
	readonly tagsSection: Locator;
	readonly sortBySection: Locator;
	/** Character portrait images in the expanded Characters section / cards. */
	readonly characterPortraits: Locator;

	constructor(page: Page) {
		this.page = page;
		this.searchBox = page.getByRole("button", { name: "Type to search..." });
		this.filterButton = page.getByRole("button", {
			name: "Filter",
			exact: true,
		});
		this.copyConfigButtons = page.getByRole("button", { name: "Copy Config" });
		this.openInViewerLinks = page.locator("a", {
			has: page.getByRole("button", { name: "Open in Viewer" }),
		});
		this.charactersSection = page.getByRole("button", { name: /Characters/ });
		this.tagsSection = page.getByRole("button", { name: /Tags/ });
		this.sortBySection = page.getByRole("button", { name: /Sort by/ });
		this.characterPortraits = page.locator('img[src^="/api/assets/avatar/"]');
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

	/** Expand the Characters section and wait for its portrait picker to render. */
	async expandCharacters(): Promise<void> {
		await this.charactersSection.click();
		await expect(this.characterPortraits.first()).toBeVisible();
	}

	/**
	 * Filter to a single character: open the search popover, type the name into
	 * its cmdk input, then pick the matching option. Dispatches an include
	 * filter, which refetches `/api/db` with the narrowed query.
	 */
	async filterByCharacter(name: string): Promise<void> {
		await this.searchBox.click();
		await this.page.getByPlaceholder("Type to search...").fill(name);
		await this.page.getByRole("option", { name }).first().click();
	}
}
