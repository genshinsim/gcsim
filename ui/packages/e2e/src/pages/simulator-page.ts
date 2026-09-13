import { expect, type Locator, type Page } from "@playwright/test";

/**
 * Drives the `/simulator` route: the config editor and the Run button.
 *
 * This encodes the app's implicit boot protocol once. The app ships no
 * data-testids, so every locator here rides an observable contract (a DOM id,
 * a Blueprint class, or the button's accessible name) documented in the README.
 */
export class SimulatorPage {
	readonly page: Page;
	/** Blueprint button, accessible name "Run". */
	readonly runButton: Locator;
	/** Ace editor container (`#config_editor`) holding the config text. */
	readonly editor: Locator;

	constructor(page: Page) {
		this.page = page;
		this.runButton = page.getByRole("button", { name: "Run" });
		this.editor = page.locator("#config_editor");
	}

	/** Navigate to the simulator and wait for React to mount into `#root`. */
	async goto(): Promise<void> {
		await this.page.goto("/simulator");
		await expect(this.page.locator("#root")).not.toBeEmpty();
		await expect(this.editor).toBeVisible();
	}

	/**
	 * Wait for wasm + workers to finish loading. While loading, the Blueprint
	 * Run button carries `bp4-loading` (a spinner); readiness is that class
	 * leaving. Console emits "aggregator loaded okay" / "loading N workers"
	 * during this window.
	 */
	async waitForReady(): Promise<void> {
		await expect(this.runButton).not.toHaveClass(/bp4-loading/, {
			timeout: 60_000,
		});
	}

	/**
	 * Replace the editor contents with `cfg`. Dispatches a native paste event on
	 * Ace's proxy textarea rather than typing key-by-key — typing would trip
	 * Ace's auto-indent and bracket matching and corrupt the config.
	 */
	async setConfig(cfg: string): Promise<void> {
		const textarea = this.editor.locator("textarea.ace_text-input");
		await textarea.focus();
		await this.page.evaluate((text) => {
			const ta = document.querySelector<HTMLTextAreaElement>(
				"#config_editor textarea.ace_text-input",
			);
			if (ta == null) {
				throw new Error("config editor textarea not found");
			}
			ta.focus();
			const data = new DataTransfer();
			data.setData("text/plain", text);
			ta.dispatchEvent(
				new ClipboardEvent("paste", {
					clipboardData: data,
					bubbles: true,
					cancelable: true,
				}),
			);
		}, cfg);
	}

	/**
	 * Wait for the config to validate: the Run button becomes enabled and no
	 * "Invalid Config" error callout is present. Validation logs "all is good".
	 */
	async waitForConfigValid(): Promise<void> {
		await expect(this.runButton).toBeEnabled({ timeout: 30_000 });
		await expect(this.page.getByText("Invalid Config")).toHaveCount(0);
	}

	/** Click Run. Navigates the app to `/web`. */
	async run(): Promise<void> {
		await this.runButton.click();
		await this.page.waitForURL(/\/web$/);
	}
}
