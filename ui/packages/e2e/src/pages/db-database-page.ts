import { expect, type Locator, type Page } from "@playwright/test";

export class DbDatabasePage {
	readonly page: Page;
	readonly searchBox: Locator;
	readonly filterButton: Locator;
	readonly copyConfigButtons: Locator;
	readonly openInViewerLinks: Locator;
	readonly charactersSection: Locator;
	readonly tagsSection: Locator;
	readonly sortBySection: Locator;
	readonly charSearch: Locator;
	/** Character portrait images in the filter picker / entry cards. */
	readonly characterPortraits: Locator;

	constructor(page: Page) {
		this.page = page;
		this.searchBox = page.getByPlaceholder(
			"Search characters, authors, notes…",
		);
		this.filterButton = page.getByRole("button", {
			name: "Filter",
			exact: true,
		});
		this.copyConfigButtons = page.getByRole("button", { name: "Copy config" });
		this.openInViewerLinks = page.getByRole("link", {
			name: "Open in viewer",
		});
		this.charactersSection = page.getByText("Characters", { exact: true });
		this.tagsSection = page.getByText("Tags", { exact: true });
		this.sortBySection = page.getByText("Sort by", { exact: true });
		this.charSearch = page.getByPlaceholder("Type to search...");
		this.characterPortraits = page.locator('img[src^="/api/assets/avatar/"]');
	}

	/** Navigate to the browse route and wait for React to mount into `#root`. */
	async goto(): Promise<void> {
		await this.page.goto("/database");
		await expect(this.page.locator("#root")).not.toBeEmpty();
	}

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

	async openFilterPanel(): Promise<void> {
		await this.filterButton.click();
		await expect(this.charSearch).toBeVisible();
	}

	/** The filter sheet shows all sections flat; wait for the portrait picker. */
	async expectCharacterPicker(): Promise<void> {
		await expect(this.characterPortraits.first()).toBeVisible();
	}

	async filterByCharacter(name: string): Promise<void> {
		await this.openFilterPanel();
		await this.charSearch.fill(name);
		await this.page.getByRole("button", { name, exact: true }).first().click();
		await this.page.keyboard.press("Escape");
	}
}
