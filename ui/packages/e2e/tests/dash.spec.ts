import { installDbRoutes, test } from "../src";

test.describe("dash home", () => {
	test("renders nav, Open Simulator, and the shared-by-others section", async ({
		app,
	}) => {
		await installDbRoutes(app.page);

		await app.dash.goto();
		await app.dash.waitForLoaded();
		await app.dash.waitForSharedByOthers();

		app.console.assertNoErrors();
	});
});
