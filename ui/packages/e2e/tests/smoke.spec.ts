import { sucroseConfig, test } from "../src";

test.describe("smoke", () => {
	test("boots, runs a sim, and renders the viewer", async ({ app }) => {
		// App mounts and wasm + workers become ready (Run leaves its spinner).
		await app.boot();

		// The provided config validates (Run enables) and the sim runs to
		// completion, navigating to the viewer.
		await app.run(sucroseConfig);

		// The viewer renders: title, tab strip, and the Results tab's cards +
		// inline charts. Structural only — no numeric result assertions.
		await app.viewer.waitForViewer();
		await app.viewer.waitForResults();

		// Nothing threw or logged an error across the whole path.
		app.console.assertNoErrors();
	});
});
