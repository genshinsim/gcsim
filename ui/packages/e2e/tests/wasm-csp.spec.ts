import { expect, sucroseConfig, test } from "../src";

// The page's CSP (index.html) doesn't allow wasm: script-src has no 'wasm-unsafe-eval'. The
// workers load from their own URLs and don't inherit it, so they compile and instantiate the
// wasm; the page itself must never compile, instantiate or even receive a compiled module.
// Firefox refuses to deserialize a WebAssembly.Module on such a page: a run that routes the
// module through the page hangs there without an error. Chromium allows it, so this spec runs
// in Firefox and WebKit too.
test("runs a sim on several workers under the page's CSP", async ({
	app,
	page,
}) => {
	await page.addInitScript(() => {
		localStorage.setItem("wasm-num-workers", "3");
		const w = window as unknown as { cspViolations: string[] };
		w.cspViolations = [];
		document.addEventListener("securitypolicyviolation", (e) => {
			w.cspViolations.push(`${e.violatedDirective}: ${e.blockedURI}`);
		});
	});

	await app.boot();
	const csp = await page
		.locator('meta[http-equiv="Content-Security-Policy"]')
		.getAttribute("content");
	expect(
		csp,
		"this spec checks running without 'wasm-unsafe-eval'; update it if the CSP now allows wasm",
	).toMatch(/script-src/);
	expect(csp).not.toMatch(/'(wasm-)?unsafe-eval'/);

	await app.run(sucroseConfig.replace("iteration=1;", "iteration=30;"));
	await app.viewer.waitForViewer();

	expect(
		await page.evaluate(
			() => (window as unknown as { cspViolations: string[] }).cspViolations,
		),
	).toEqual([]);
});
