import { ENKA_UID, expect, goodImport, installEnkaRoutes, test } from "../src";

test.describe("toolbox imports", () => {
	test("GOOD import validates a pasted fixture", async ({ app }) => {
		// The toolbox renders on mount — no wasm/workers needed for import.
		await app.simulator.goto();

		const dialog = await app.simulator.openImportDialog("GO");
		await dialog.locator("textarea").fill(goodImport);
		await expect(dialog.getByText("Data parsed successfully.")).toBeVisible();

		// Loading the parsed fixture toasts success and closes the dialog.
		await dialog.getByRole("button", { name: "Import" }).click();
		await expect(
			app.page.getByText("GOOD data imported successfully"),
		).toBeVisible();
		await expect(dialog).toBeHidden();

		app.console.assertNoErrors();
	});

	test("Enka import validates a stubbed fixture", async ({ app }) => {
		await installEnkaRoutes(app.page);
		await app.simulator.goto();

		const dialog = await app.simulator.openImportDialog("Enka");
		await dialog.getByPlaceholder("Paste UID here").fill(ENKA_UID);
		await dialog.getByRole("button", { name: "Import" }).click();

		await expect(
			dialog.getByText("Data retrieved successfully.", { exact: false }),
		).toBeVisible();
		await expect(dialog.getByText("bennett")).toBeVisible();

		app.console.assertNoErrors();
	});
});
