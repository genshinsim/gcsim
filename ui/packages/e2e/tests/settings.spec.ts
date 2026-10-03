import { expect, test } from "../src";

test.describe("settings page", () => {
	test("keeps URL edits in server mode", async ({ app, page }) => {
		await app.simulator.goto();
		await page.getByRole("link", { name: "Settings", exact: true }).click();
		await expect(page).toHaveURL(/\/settings$/);

		const nav = page.getByRole("navigation");
		await expect(nav.getByText("WASM", { exact: true })).toBeVisible();
		await page.locator("#server-mode-switch").click();
		await expect(nav.getByText("Server", { exact: true })).toBeVisible();

		const url = page.locator("#server-mode-url");
		await expect(url).toBeVisible();

		await url.fill("http://127.0.0.1:1234");
		await expect(url).toHaveValue("http://127.0.0.1:1234");
	});

	test("applies the saved theme without the bundle", async ({ page }) => {
		await page.goto("/settings");
		await page
			.locator("label", {
				has: page.getByRole("radio", { name: /Ember\s*Light/ }),
			})
			.click();
		await expect(page.locator("html")).toHaveAttribute("data-theme", "ember-l");

		await page.route("**/src/index.tsx", (route) => route.abort());
		await page.reload();
		await expect(page.locator("html")).toHaveAttribute("data-theme", "ember-l");
	});
});
