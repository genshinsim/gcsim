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
		this.heading = page.getByText(/Showing entries with the same team for id:/);
		this.mainCard = this.heading.locator("..");
		this.existingSection = page
			.getByText("Existing sims with same characters")
			.locator("..");
	}

	async goto(id: string): Promise<void> {
		await this.page.goto(`/id/${id}`);
		await expect(this.page.locator("#root")).not.toBeEmpty();
		await expect(this.heading).toBeVisible();
	}

	async waitForEntry(
		chars: readonly string[],
		sourceTag: string,
	): Promise<void> {
		for (const name of chars) {
			await expect(this.mainCard.locator(`img[alt="${name}"]`)).toBeVisible();
		}
		for (const chip of [
			"mode",
			"target count",
			"dps/target",
			"avg sim time",
			"created",
			sourceTag,
		]) {
			await expect(
				this.mainCard.getByText(chip, { exact: false }),
			).toBeVisible();
		}
	}

	async waitForControls(): Promise<void> {
		await expect(
			this.mainCard.getByRole("button", { name: "Copy Reject" }),
		).toBeVisible();
		await expect(
			this.mainCard.getByRole("button", { name: "Copy Approve" }),
		).toBeVisible();
		await expect(
			this.mainCard.getByRole("button", { name: "Result Viewer" }),
		).toBeVisible();
	}

	async waitForExistingSims(): Promise<void> {
		await expect(this.existingSection).toBeVisible();
		const replace = this.existingSection.getByRole("button", {
			name: "Replace This",
		});
		const empty = this.existingSection.getByText("Nothing found");
		await expect(replace.first().or(empty)).toBeVisible();
	}

	/** Click a moderation copy button and return the resulting clipboard text. */
	async copyCommand(name: "Copy Reject" | "Copy Approve"): Promise<string> {
		await this.mainCard.getByRole("button", { name }).click();
		return this.page.evaluate(() => navigator.clipboard.readText());
	}
}
