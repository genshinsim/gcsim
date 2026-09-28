import { expect, installDbRoutes, sucroseConfig, test } from "../src";

test.describe("routes", () => {
	const pages: [string, string][] = [
		["/simulator", "gcsim - simulator"],
		["/web", "gcsim - viewer"],
		["/sample/upload", "gcsim - sample"],
		["/account", "gcsim - account"],
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

	test("nav links move between pages", async ({ page }) => {
		await installDbRoutes(page);
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
