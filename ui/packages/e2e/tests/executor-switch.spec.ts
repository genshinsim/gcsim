import type { Page } from "@playwright/test";
import { expect, sucroseConfig, test } from "../src";

const SERVER = "http://127.0.0.1:54321";

async function toggleExecutor(page: Page, chip: "WASM" | "Server") {
	await page
		.getByRole("navigation")
		.getByRole("button", { name: chip, exact: true })
		.click();
	await page.locator("#server-mode-switch").click();
	await page.keyboard.press("Escape");
	await expect(page.getByRole("dialog")).toHaveCount(0);
}

test.describe("executor switch", () => {
	test("server → WASM runs on WASM without a reload", async ({ app, page }) => {
		await page.route(`${SERVER}/**`, (route) => route.abort());
		await page.addInitScript(() =>
			localStorage.setItem("use-server-mode", "true"),
		);
		await app.simulator.goto();

		await toggleExecutor(page, "Server");
		await app.simulator.waitForReady();
		await app.run(sucroseConfig);
		await app.viewer.waitForViewer();
	});

	test("WASM → server sends to the server without a reload", async ({
		app,
		page,
	}) => {
		await app.boot();
		const ready = page.waitForRequest(`${SERVER}/ready/**`);
		await page.route(`${SERVER}/**`, (route) => route.abort());

		await toggleExecutor(page, "WASM");
		await ready;
	});
});
