import { sucroseConfig, test } from "../src";

test.describe("viewer tabs", () => {
	test("Config and Sample tabs render content after a run", async ({ app }) => {
		await app.boot();
		await app.run(sucroseConfig);
		await app.viewer.waitForViewer();
		await app.viewer.waitForResults();

		await app.viewer.openConfig();
		await app.viewer.openSample();

		app.console.assertNoErrors();
	});
});
