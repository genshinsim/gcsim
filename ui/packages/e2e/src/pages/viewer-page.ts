import { expect, type Locator, type Page } from "@playwright/test";

/**
 * Drives the `/web` viewer route: the tab strip and the Results tab content.
 *
 * Locators ride observable contracts (page title, tab accessible names,
 * Blueprint `bp4-card`, inline `<svg>`) since the app ships no data-testids.
 */
export class ViewerPage {
	readonly page: Page;
	readonly resultsTab: Locator;
	readonly configTab: Locator;
	readonly sampleTab: Locator;
	/** Blueprint cards on the Results tab. */
	readonly cards: Locator;
	/** Inline visx charts on the Results tab. */
	readonly charts: Locator;

	constructor(page: Page) {
		this.page = page;
		this.resultsTab = page.getByRole("tab", { name: "Results" });
		this.configTab = page.getByRole("tab", { name: "Config" });
		this.sampleTab = page.getByRole("tab", { name: "Sample" });
		this.cards = page.locator(".bp4-card");
		this.charts = page.locator(".bp4-card svg");
	}

	/** Assert the viewer chrome rendered: title and the three-tab strip. */
	async waitForViewer(): Promise<void> {
		await expect(this.page).toHaveTitle("gcsim - viewer");
		await expect(this.resultsTab).toBeVisible();
		await expect(this.configTab).toBeVisible();
		await expect(this.sampleTab).toBeVisible();
	}

	/**
	 * Assert the Results tab rendered the sim output: at least one card and at
	 * least one inline chart. Structural only — never asserts result numbers.
	 */
	async waitForResults(): Promise<void> {
		await expect(this.cards.first()).toBeVisible({ timeout: 30_000 });
		await expect(this.charts.first()).toBeVisible({ timeout: 30_000 });
	}
}
