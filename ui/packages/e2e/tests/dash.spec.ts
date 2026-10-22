import { expect, installOfflineRoutes, test } from "../src";

const layouts = [
	{ name: "desktop", viewport: { width: 1280, height: 720 } },
	{ name: "mobile", viewport: { width: 390, height: 844 } },
];

test.describe("dash home", () => {
	for (const { name, viewport } of layouts) {
		test(`${name}: renders without the KQM DB section`, async ({ app }) => {
			await app.page.setViewportSize(viewport);
			await installOfflineRoutes(app.page);
			const kqmRequests: string[] = [];
			app.page.on("request", (req) => {
				if (new URL(req.url()).host === "db.kqm.gg") {
					kqmRequests.push(req.url());
				}
			});

			await app.dash.goto();
			if (name === "desktop") {
				await app.dash.waitForLoaded();
			} else {
				// The mobile nav collapses its Simulator link behind the menu.
				await expect(app.dash.openSimulator).toBeVisible();
			}

			await expect(
				app.page.getByRole("heading", { name: /KQM DB/ }),
			).toHaveCount(0);
			await expect(
				app.page.locator("a[href^='https://db.kqm.gg']"),
			).toHaveCount(0);
			expect(kqmRequests).toEqual([]);
			app.console.assertNoErrors();
		});
	}
});
