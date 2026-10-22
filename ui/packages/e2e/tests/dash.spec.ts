import { type DashPage, expect, installOfflineRoutes, test } from "../src";

const layouts = [
	{
		name: "desktop",
		viewport: { width: 1280, height: 720 },
		waitForLoaded: (dash: DashPage) => dash.waitForLoaded(),
	},
	{
		name: "mobile",
		viewport: { width: 390, height: 844 },
		waitForLoaded: (dash: DashPage) => dash.waitForLoadedMobile(),
	},
];

test.describe("dash home", () => {
	for (const { name, viewport, waitForLoaded } of layouts) {
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
			await waitForLoaded(app.dash);

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
