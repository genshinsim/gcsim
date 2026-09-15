import { expect, test } from "../../src";

const PAGES: {
	section: string;
	path: string;
	title: string;
	heading: string;
	images: boolean;
}[] = [
	{
		section: "Introduction",
		path: "/",
		title: "Introduction",
		heading: "Introduction",
		images: true,
	},
	{
		section: "Guides",
		path: "/guides/building_a_simulation_basic_tutorial",
		title: "Building a Simulation from Scratch",
		heading: "Building a Simulation from Scratch",
		images: true,
	},
	{
		section: "Reference",
		path: "/reference/config",
		title: "Config File",
		heading: "Config File",
		images: false,
	},
	{
		section: "Game Mechanics",
		path: "/mechanics/aoe",
		title: "Area of Effect",
		heading: "Area of Effect",
		images: true,
	},
];

const SECTIONS = PAGES.map((p) => p.section);

function routeRegex(path: string): RegExp {
	const escaped = path.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
	return new RegExp(`${escaped}/?$`);
}

test.describe("docs smoke", () => {
	test("sidebar lists every top-level section", async ({ docs }) => {
		await docs.goto("/");
		for (const label of SECTIONS) {
			await expect(docs.sidebarLink(label)).toBeVisible();
		}
		docs.console.assertNoCrashes();
	});

	for (const page of PAGES) {
		test(`${page.section} (${page.path}) renders`, async ({ docs }) => {
			await docs.goto(page.path);

			await expect(docs.page).toHaveURL(routeRegex(page.path));
			await expect(docs.page).toHaveTitle(`${page.title} | gcsim Docs`);
			await expect(docs.heading(page.heading)).toBeVisible();
			await expect(docs.content).not.toBeEmpty();

			if (page.images) {
				expect(await docs.content.locator("img").count()).toBeGreaterThan(0);
			}
			await docs.assertImagesLoaded(docs.content);

			docs.console.assertNoCrashes();
		});
	}

	test("an uncaught page exception is recorded as a crash", async ({
		docs,
	}) => {
		await docs.goto("/");
		await docs.page.evaluate(() => {
			setTimeout(() => {
				throw new Error("gcsim-e2e synthetic crash");
			}, 0);
		});
		await expect.poll(() => docs.console.crashes.length).toBeGreaterThan(0);
		expect(docs.console.crashes.join("\n")).toContain(
			"gcsim-e2e synthetic crash",
		);
		expect(() => docs.console.assertNoCrashes()).toThrow();
	});

	test("a console error is not a crash", async ({ docs }) => {
		await docs.goto("/");
		await docs.page.evaluate(() =>
			console.error("gcsim-e2e synthetic 404-like error"),
		);
		await expect.poll(() => docs.console.errors.length).toBeGreaterThan(0);
		expect(docs.console.crashes).toEqual([]);
	});
});
