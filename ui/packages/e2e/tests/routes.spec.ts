import { expect, installOfflineRoutes, sucroseConfig, test } from "../src";

test.describe("routes", () => {
	const pages: [string, string][] = [
		["/simulator", "gcsim - simulator"],
		["/web", "gcsim - viewer"],
		["/sample/upload", "gcsim - sample"],
		["/no-such-page", "gcsim - simulation impact"],
	];
	for (const [path, title] of pages) {
		test(`${path} renders its page`, async ({ page }) => {
			await page.goto(path);
			await expect(page).toHaveTitle(title);
		});
	}

	const redirects: [string, string][] = [
		["/v3/viewer/share/abc", "/sh/abc"],
		["/viewer/share/abc", "/sh/abc"],
		["/s/abc", "/sh/abc"],
		["/viewer/web", "/web"],
		["/viewer/local", "/local"],
		["/simple", "/simulator"],
		["/advanced", "/simulator"],
		["/viewer", "/simulator"],
	];
	for (const [from, to] of redirects) {
		test(`${from} redirects to ${to}`, async ({ page }) => {
			await page.route("**/api/share/**", (route) =>
				route.fulfill({ status: 404, body: "" }),
			);
			await page.goto(from);
			await expect(page).toHaveURL(new RegExp(`${to}$`));
		});
	}

	test("/db stays put and shows not-found", async ({ page }) => {
		await page.goto("/db");
		await expect(page).toHaveURL(/\/db$/);
		await expect(
			page.getByText("This page is not implemented yet. Stay tuned!"),
		).toBeVisible();
	});

	test("/db/:id shows the unavailable message without fetching", async ({
		page,
	}) => {
		const shareRequests: string[] = [];
		page.on("request", (req) => {
			if (new URL(req.url()).pathname.startsWith("/api/share/db")) {
				shareRequests.push(req.url());
			}
		});
		await page.goto("/db/abc");
		await expect(
			page.getByText("This database entry is no longer available on gcsim."),
		).toBeVisible();
		await expect(page.getByRole("tab", { name: "Results" })).toHaveCount(0);
		expect(shareRequests).toEqual([]);
	});

	test("nav links move between pages", async ({ page }) => {
		await installOfflineRoutes(page);
		await page.goto("/");
		await page.locator("nav a[href='/simulator']").first().click();
		await expect(page).toHaveURL(/\/simulator$/);
		await expect(page).toHaveTitle("gcsim - simulator");
		await page.locator("nav a[href='/']").click();
		await expect(page).toHaveURL(/\/$/);
		await expect(page).toHaveTitle("gcsim - simulation impact");
	});

	test("run goes to /web and send to simulator comes back", async ({
		app,
		page,
	}) => {
		await app.boot();
		await app.run(sucroseConfig);
		await expect(page).toHaveURL(/\/web$/);
		await app.viewer.waitForViewer();
		await app.viewer.waitForResults();

		await page.getByRole("button", { name: "Send to Simulator" }).click();
		await page.getByRole("button", { name: "Continue" }).click();
		await expect(page).toHaveURL(/\/simulator$/);
		await expect(page).toHaveTitle("gcsim - simulator");

		app.console.assertNoErrors();
	});
});
