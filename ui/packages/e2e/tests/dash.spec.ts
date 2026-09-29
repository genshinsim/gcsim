import { expect, installOfflineRoutes, kqmEntries, test } from "../src";

test.describe("dash home", () => {
	test("renders nav, Open Simulator, and KQM DB's shared-by-others cards", async ({
		app,
	}) => {
		await installOfflineRoutes(app.page);
		const gcsimDbRequests: string[] = [];
		app.page.on("request", (req) => {
			const url = new URL(req.url());
			if (url.host !== "db.kqm.gg" && url.pathname.startsWith("/api/db")) {
				gcsimDbRequests.push(req.url());
			}
		});

		await app.dash.goto();
		await app.dash.waitForLoaded();
		await app.dash.waitForSharedByOthers();

		await expect(app.dash.teamCards).toHaveCount(kqmEntries.length);
		await expect(app.dash.teamCards.first()).toHaveAttribute(
			"href",
			`https://db.kqm.gg/db/${kqmEntries[0]._id}`,
		);
		await expect(app.dash.teamCards.first()).toContainText("mode TTK");
		await expect(app.dash.teamCards.nth(1)).toContainText("mode duration");
		await expect(
			app.page.getByRole("link", { name: /View all/ }).first(),
		).toHaveAttribute("href", "https://db.kqm.gg");
		await expect(app.page.getByRole("link", { name: "Teams DB" })).toHaveCount(
			0,
		);
		expect(gcsimDbRequests).toEqual([]);
		app.console.assertNoErrors();
	});
});
