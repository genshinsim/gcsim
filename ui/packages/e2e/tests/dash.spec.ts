import { installDbRoutes, test } from "../src";

test.describe("dash home", () => {
	test("renders nav, Open Simulator, and the shared-by-others section", async ({
		app,
	}) => {
		// The shared-by-others cards fetch `/api/db`; stub it (plus the GitHub release
		// and assets) so the home page renders offline and deterministically.
		// `installDbRoutes` provides exactly that offline stub set.
		await installDbRoutes(app.page);

		await app.dash.goto();
		await app.dash.waitForLoaded();
		await app.dash.waitForSharedByOthers();

		app.console.assertNoErrors();
	});
});
