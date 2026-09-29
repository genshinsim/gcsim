import { expect, installOfflineRoutes, test } from "../src";

test.describe("dash home", () => {
	test("renders nav and Open Simulator without fetching the db", async ({
		app,
	}) => {
		await installOfflineRoutes(app.page);
		const dbRequests: string[] = [];
		app.page.on("request", (req) => {
			if (new URL(req.url()).pathname.startsWith("/api/db")) {
				dbRequests.push(req.url());
			}
		});

		await app.dash.goto();
		await app.dash.waitForLoaded();

		await expect(
			app.page.getByRole("heading", { name: "Shared by others" }),
		).toHaveCount(0);
		await expect(app.page.getByRole("link", { name: "Teams DB" })).toHaveCount(
			0,
		);
		expect(dbRequests).toEqual([]);
		app.console.assertNoErrors();
	});
});
