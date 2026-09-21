import { expect, type Locator, type Page } from "@playwright/test";
import { ConsoleMonitor } from "./console-monitor";

export class TaghelperHarness {
	readonly page: Page;
	readonly console: ConsoleMonitor;
	readonly heading: Locator;
	readonly mainCard: Locator;
	readonly existingSection: Locator;

	constructor(page: Page) {
		this.page = page;
		this.console = new ConsoleMonitor(page);
		// The "Under review" section holds the submission card and its actions.
		this.heading = page.getByText("Under review");
		this.mainCard = page.locator("section", { has: this.heading });
		this.existingSection = page.locator("section", {
			has: page.getByText("Existing sims with the same team"),
		});
	}

	async goto(id: string): Promise<void> {
		await this.page.goto(`/id/${id}`);
		await expect(this.page.locator("#root")).not.toBeEmpty();
		await expect(this.heading).toBeVisible();
	}

	async waitForEntry(chars: readonly string[]): Promise<void> {
		for (const name of chars) {
			await expect(
				this.mainCard.locator(`img[alt="${name}"]`).first(),
			).toBeVisible();
		}
		// Summary meta chips: sim mode and mean sim duration.
		await expect(this.mainCard.getByText("mode TTK")).toBeVisible();
		await expect(this.mainCard.getByText("90.0s")).toBeVisible();
	}

	async waitForControls(): Promise<void> {
		await expect(
			this.mainCard.getByRole("button", { name: "Copy reject" }),
		).toBeVisible();
		await expect(
			this.mainCard.getByRole("button", { name: "Copy approve" }),
		).toBeVisible();
		await expect(
			this.mainCard.getByRole("link", { name: "Result viewer" }),
		).toBeVisible();
	}

	async waitForExistingSims(): Promise<void> {
		await expect(this.existingSection).toBeVisible();
		const replace = this.existingSection.getByRole("button", {
			name: "Replace",
		});
		const empty = this.existingSection.getByText(
			"No existing sims share this team.",
		);
		await expect(replace.first().or(empty)).toBeVisible();
	}

	/** Click a moderation copy button and return the resulting clipboard text. */
	async copyCommand(name: "Copy reject" | "Copy approve"): Promise<string> {
		await this.mainCard.getByRole("button", { name }).click();
		return this.page.evaluate(() => navigator.clipboard.readText());
	}
}
