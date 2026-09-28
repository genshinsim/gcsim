import { readFileSync } from "node:fs";
import { expect, test } from "../src";

const sample = readFileSync(
	new URL(
		"../../components/src/SampleLog/testdata/sample.json",
		import.meta.url,
	),
	"utf8",
);

test("/sample/local loads the sample from the local server", async ({
	page,
}) => {
	await page.route("http://127.0.0.1:8381/sample", (route) =>
		route.fulfill({
			status: 200,
			contentType: "application/json",
			body: sample,
		}),
	);
	await page.goto("/sample/local");
	await expect(page).toHaveTitle("gcsim - local sample");
	await expect(page.getByRole("button", { name: "Download" })).toBeVisible();
});
