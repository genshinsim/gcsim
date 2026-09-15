import { expect, MAIN_ID, mainEntry, SOURCE_TAG_NAME, test } from "../../src";

// The copy assertion reads the clipboard, which needs the permission granted.
test.use({ permissions: ["clipboard-read", "clipboard-write"] });

const chars = mainEntry.summary.char_names;

test.describe("taghelper smoke", () => {
	test("renders the moderation view for an entry", async ({ taghelper }) => {
		// The `/id/:id` view boots and renders the main entry.
		await taghelper.goto(MAIN_ID);

		// Main card: team portraits and the summary stat chips.
		await taghelper.waitForEntry(chars, SOURCE_TAG_NAME);

		// Moderation controls and the "existing sims" section render.
		await taghelper.waitForControls();
		await taghelper.waitForExistingSims();

		// Copy Approve writes the expected slash command to the clipboard.
		expect(await taghelper.copyCommand("Copy Approve")).toBe(
			`/approve id:${MAIN_ID}`,
		);

		// Nothing threw or logged an error across the whole path.
		taghelper.console.assertNoErrors();
	});
});
