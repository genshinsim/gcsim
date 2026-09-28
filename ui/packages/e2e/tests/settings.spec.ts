import { expect, test } from "../src";

test.describe("settings dialog", () => {
	test("keeps URL edits in server mode", async ({ app, page }) => {
		await app.simulator.goto();
		const settings = page.getByRole("button", { name: "Settings" });
		await settings.click();
		// switching modes remounts the app, which closes the dialog
		await page.locator("#server-mode-switch").click();
		await settings.click();

		const url = page.locator("#server-mode-url");
		await expect(url).toBeVisible();

		await url.fill("http://127.0.0.1:1234");
		await expect(url).toHaveValue("http://127.0.0.1:1234");
	});
});
