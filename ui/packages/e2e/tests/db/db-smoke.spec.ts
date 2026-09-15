import { expect, test } from "../../src";

test.describe("db smoke", () => {
	test("home page renders", async ({ db }) => {
		await db.home.goto();

		// Welcome copy and the "Get started" CTA render.
		await db.home.waitForLoaded();

		// Nothing threw or logged an error.
		db.console.assertNoErrors();
	});

	test("browse, filter panel, and character filter", async ({ db }) => {
		await db.database.goto();

		// The browse view renders: result count, search box, filter funnel, and
		// at least one entry card with its Copy Config / Open in Viewer controls.
		await db.database.waitForBrowse();

		// Open in Viewer targets the web app's cross-app /db/:id route. Assert the
		// target only — following it is a separate concern.
		await expect(db.database.openInViewerLinks.first()).toHaveAttribute(
			"href",
			/\/db\/[a-z0-9]+$/,
		);

		// The filter drawer exposes the Characters, Tags, and Sort by sections.
		await db.database.openFilterPanel();
		await expect(db.database.charactersSection).toBeVisible();
		await expect(db.database.tagsSection).toBeVisible();
		await expect(db.database.sortBySection).toBeVisible();

		// Exercising a character filter narrows the list (2 -> 1 simulations).
		await db.page.keyboard.press("Escape");
		await db.database.filterByCharacter("Nahida");
		await db.database.expectShowing(1);

		// Nothing threw or logged an error across the whole path.
		db.console.assertNoErrors();
	});
});
