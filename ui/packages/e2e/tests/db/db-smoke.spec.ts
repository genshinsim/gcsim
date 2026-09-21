import { expect, test } from "../../src";

test.describe("db smoke", () => {
	test("home page renders", async ({ db }) => {
		await db.home.goto();

		// The hero, "Browse database" CTA, and "What's new" section render.
		await db.home.waitForLoaded();

		// Nothing threw or logged an error.
		db.console.assertNoErrors();
	});

	test("browse, filter sheet, and character filter", async ({ db }) => {
		await db.database.goto();

		// The browse view renders: result count, search box, filter trigger, and
		// at least one entry card with its Copy config / Open in viewer controls.
		await db.database.waitForBrowse();

		// Open in viewer targets the web app's cross-app /db/:id route. Assert the
		// target only — following it is a separate concern.
		await expect(db.database.openInViewerLinks.first()).toHaveAttribute(
			"href",
			/\/db\/[a-z0-9]+$/,
		);

		// The filter sheet exposes the Characters, Tags, and Sort by sections.
		await db.database.openFilterPanel();
		await expect(db.database.charactersSection).toBeVisible();
		await expect(db.database.tagsSection).toBeVisible();
		await expect(db.database.sortBySection).toBeVisible();

		// The Characters section expands to its portrait picker.
		await db.database.expandCharacters();
		await db.page.keyboard.press("Escape");

		// Exercising a character filter narrows the list (2 -> 1 simulations).
		await db.database.filterByCharacter("Nahida");
		await db.database.expectShowing(1);

		// Nothing threw or logged an error across the whole path.
		db.console.assertNoErrors();
	});
});
