import type { Page, Route } from "@playwright/test";
import { type AppHarness, expect, sucroseConfig, test } from "../src";

const SHARE_URL = "**/api/share/abc";
const AUTH = "auth-hash-1";

let cachedResult: string | undefined;

async function simResult(app: AppHarness, page: Page): Promise<string> {
	if (cachedResult == null) {
		await app.boot();
		await app.run(sucroseConfig);
		cachedResult = (await leaveToSaveLastRun(page)) ?? undefined;
	}
	expect(cachedResult).toBeTruthy();
	return cachedResult as string;
}

async function leaveToSaveLastRun(page: Page): Promise<string | null> {
	await page.goto("/account");
	return page.evaluate(() => localStorage.getItem("redux-local-results"));
}

function fulfillShare(body: string) {
	return (route: Route) =>
		route.fulfill({
			status: 200,
			contentType: "application/json",
			headers: {
				"x-gcsim-share-auth": AUTH,
				"access-control-allow-origin": "*",
				"access-control-expose-headers": "x-gcsim-share-auth",
			},
			body,
		});
}

test.describe("share viewer", () => {
	test("?tab=sample opens the Sample tab", async ({ app, page }) => {
		await page.route(SHARE_URL, fulfillShare(await simResult(app, page)));
		await page.goto("/sh/abc?tab=sample");
		await expect(page).toHaveTitle("gcsim sh - abc");
		await expect(app.viewer.sampleTab).toHaveAttribute("aria-selected", "true");
		await expect(app.viewer.generateButton).toBeVisible();
	});

	test("#tab=sample&sample= redirects to search params", async ({
		app,
		page,
	}) => {
		await page.route(SHARE_URL, fulfillShare(await simResult(app, page)));
		await page.goto("/sh/abc#tab=sample&sample=123");
		await expect(page).toHaveURL(/\/sh\/abc\?tab=sample&seed=123$/);
		await expect(app.viewer.sampleTab).toHaveAttribute("aria-selected", "true");
	});

	test("an old share link keeps its #tab= through both redirects", async ({
		app,
		page,
	}) => {
		await page.route(SHARE_URL, fulfillShare(await simResult(app, page)));
		await page.goto("/viewer/share/abc#tab=config");
		await expect(page).toHaveURL(/\/sh\/abc\?tab=config$/);
	});

	test("tab clicks move the tab search param without refetching", async ({
		app,
		page,
	}) => {
		const fulfill = fulfillShare(await simResult(app, page));
		let fetches = 0;
		await page.route(SHARE_URL, (route) => {
			fetches++;
			return fulfill(route);
		});
		await page.goto("/sh/abc");
		await app.viewer.openConfig();
		await expect(page).toHaveURL(/\/sh\/abc\?tab=config$/);
		await app.viewer.openSample();
		await page.goBack();
		await page.goBack();
		await expect(app.viewer.resultsTab).toHaveAttribute(
			"aria-selected",
			"true",
		);
		expect(fetches).toBe(1);
	});

	test("a failing load shows the error dialog and retry reloads", async ({
		app,
		page,
	}) => {
		const body = await simResult(app, page);
		let fail = true;
		await page.route(SHARE_URL, (route) =>
			fail
				? route.fulfill({ status: 500, body: "" })
				: fulfillShare(body)(route),
		);
		await page.goto("/sh/abc");
		const retry = page.getByRole("button", { name: "Retry" });
		await expect(retry).toBeVisible();

		fail = false;
		await retry.click();
		await expect(retry).toBeHidden();
		await app.viewer.waitForResults();
	});

	test("sharing a local result sends its loaded auth hash", async ({
		app,
		page,
	}) => {
		await page.route(
			"http://127.0.0.1:8381/data",
			fulfillShare(await simResult(app, page)),
		);
		const posted = new Promise<string | undefined>((resolve) =>
			page.route("**/api/share", (route) => {
				resolve(route.request().headers()["x-gcsim-share-auth"]);
				return route.fulfill({ status: 200, body: JSON.stringify("xyz") });
			}),
		);
		await page.goto("/local");
		await app.viewer.waitForResults();
		await page.getByRole("button", { name: "Share" }).click();
		expect(await posted).toBe(AUTH);
	});
});
