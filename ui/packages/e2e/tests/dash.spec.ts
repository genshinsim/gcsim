import { installDbRoutes, test } from "../src";

test.describe("dash home", () => {
	test("renders nav, Get started, and the featured submission", async ({
		app,
	}) => {
		// The featured card fetches `/api/db`; stub it (plus the GitHub release
		// and assets) so the home page renders offline and deterministically.
		// `installDbRoutes` provides exactly that offline stub set.
		await installDbRoutes(app.page);

		await app.dash.goto();
		await app.dash.waitForLoaded();
		await app.dash.waitForFeatured();

		app.console.assertNoErrors();
	});
});
