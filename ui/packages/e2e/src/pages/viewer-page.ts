import { expect, type Locator, type Page } from "@playwright/test";

/**
 * Drives the `/web` viewer route: the tab strip and the Results tab content.
 *
 * Locators ride observable contracts (page title, tab accessible names, the
 * `data-slot="card"` hook on ported cards, chart `role="img"`) since the app
 * ships no data-testids.
 */
export class ViewerPage {
	readonly page: Page;
	readonly resultsTab: Locator;
	readonly configTab: Locator;
	readonly sampleTab: Locator;
	/** Ported cards on the Results tab (shadcn `data-slot="card"`). */
	readonly cards: Locator;
	/** Inline visx charts on the Results tab (each svg is `role="img"`). */
	readonly charts: Locator;
	/** Ace editor container (`#config_editor`) shown on the Config tab. */
	readonly configEditor: Locator;
	/** The Sample tab's "Generate" button (shown before a sample is generated). */
	readonly generateButton: Locator;

	constructor(page: Page) {
		this.page = page;
		this.resultsTab = page.getByRole("tab", { name: "Results" });
		this.configTab = page.getByRole("tab", { name: "Config" });
		this.sampleTab = page.getByRole("tab", { name: "Sample" });
		this.cards = page.locator('[data-slot="card"]');
		this.charts = page.locator('svg[role="img"]');
		this.configEditor = page.locator("#config_editor");
		this.generateButton = page.getByRole("button", { name: "Generate" });
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

	/**
	 * Open the Config tab and assert it rendered the config editor holding the
	 * run's config text. Structural only.
	 */
	async openConfig(): Promise<void> {
		await this.configTab.click();
		await expect(this.configEditor).toBeVisible();
		await expect(this.configEditor).not.toBeEmpty();
	}

	/**
	 * Open the Sample tab and assert it rendered its content — the "Generate"
	 * control shown before a sample is generated. Structural only; does not
	 * generate a sample.
	 */
	async openSample(): Promise<void> {
		await this.sampleTab.click();
		await expect(this.generateButton).toBeVisible();
	}
}
