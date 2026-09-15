import { expect, test } from "../../src";

const SIDEBAR = [
	"Introduction",
	"Getting Started",
	"Guides",
	"Reference",
	"Game Mechanics",
];

test.describe("docs smoke", () => {
	test("landing renders and navigates into a nested doc", async ({ docs }) => {
		// Landing page (`/`) renders: title and doc body.
		await docs.goto("/");
		await expect(docs.page).toHaveTitle("Introduction | gcsim Docs");
		await expect(docs.content).not.toBeEmpty();

		// The sidebar lists every top-level section.
		for (const label of SIDEBAR) {
			await expect(docs.sidebarLink(label)).toBeVisible();
		}

		// A nested doc navigates and renders its own title, heading, and body.
		await docs.sidebarLink("Getting Started").click();
		await docs.page.waitForURL(/\/get-started\/?$/);
		await expect(docs.page).toHaveTitle("Getting Started | gcsim Docs");
		await expect(docs.heading("Getting Started")).toBeVisible();
		await expect(docs.content).not.toBeEmpty();

		// Nothing threw or logged an error across the whole path.
		docs.console.assertNoErrors();
	});
});
