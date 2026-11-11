import type { Locator, Page } from "@playwright/test";
import { expect, SimulatorPage, sucroseConfig, test } from "../src";

const longConfig = `${sucroseConfig}\n${Array.from(
	{ length: 50 },
	(_, i) => `# filler line ${i + 1}`,
).join("\n")}\n`;

async function openWith(
	page: Page,
	size: { width: number; height: number },
	cfg: string,
): Promise<SimulatorPage> {
	await page.setViewportSize(size);
	const sim = new SimulatorPage(page);
	await sim.goto();
	await sim.setConfig(cfg);
	return sim;
}

async function box(locator: Locator) {
	const b = await locator.boundingBox();
	if (!b) throw new Error("element not rendered");
	return { top: b.y, bottom: b.y + b.height };
}

const scrollerOf = (sim: SimulatorPage) => sim.editor.locator(".cm-scroller");

const overflow = (locator: Locator) =>
	locator.evaluate((el) => el.scrollHeight - el.clientHeight);

const scrollContent = (sim: SimulatorPage, to: "middle" | "top") =>
	sim.content.evaluate((el, to) => {
		el.scrollTop = to === "top" ? 0 : (el.scrollHeight - el.clientHeight) / 2;
	}, to);

test.describe("simulator layout", () => {
	for (const size of [
		{ width: 1440, height: 700 },
		{ width: 1440, height: 900 },
		{ width: 1440, height: 1080 },
		{ width: 1440, height: 1440 },
		{ width: 390, height: 844 },
	]) {
		test(`a long config scrolls only the page at ${size.width}x${size.height}`, async ({
			page,
		}) => {
			const sim = await openWith(page, size, longConfig);
			expect(await overflow(scrollerOf(sim))).toBeLessThanOrEqual(0);
			expect(await overflow(sim.content)).toBeGreaterThan(0);
		});
	}

	test("a short config fills the space above the Run bar without scrolling", async ({
		page,
	}) => {
		const sim = await openWith(
			page,
			{ width: 1440, height: 1440 },
			"active sucrose;",
		);
		const editor = await box(sim.editor);
		const runBar = await box(sim.runBar);
		const content = await box(sim.content);
		expect(await overflow(sim.content)).toBeLessThanOrEqual(0);
		expect(runBar.top - editor.bottom).toBeGreaterThanOrEqual(0);
		expect(runBar.top - editor.bottom).toBeLessThan(8);
		expect(content.bottom - runBar.bottom).toBeLessThan(16);
	});

	test("a short viewport keeps the editor at least 10 lines tall", async ({
		page,
	}) => {
		const sim = await openWith(
			page,
			{ width: 1440, height: 400 },
			"active sucrose;",
		);
		const line = await box(sim.editor.locator(".cm-line").first());
		const scroller = await box(scrollerOf(sim));
		expect(scroller.bottom - scroller.top).toBeGreaterThanOrEqual(
			(line.bottom - line.top) * 10,
		);
		expect(await overflow(sim.content)).toBeGreaterThan(0);
	});

	test("the footer shows everywhere but the simulator", async ({ page }) => {
		await new SimulatorPage(page).goto();
		await expect(page.getByRole("contentinfo")).toHaveCount(0);
		for (const path of ["/", "/web"]) {
			await page.goto(path);
			await expect(page.getByRole("contentinfo")).toBeVisible();
		}
	});

	test("the toolbar and Run bar stay pinned while scrolling", async ({
		page,
	}) => {
		const sim = await openWith(page, { width: 1440, height: 700 }, longConfig);
		await scrollContent(sim, "middle");
		const content = await box(sim.content);
		expect((await box(sim.toolbar)).top).toBeCloseTo(content.top, 0);
		expect((await box(sim.runBar)).bottom).toBeCloseTo(content.bottom, 0);
	});

	test("the cursor stays clear of both pinned bars", async ({ page }) => {
		const sim = await openWith(page, { width: 1440, height: 700 }, longConfig);
		await scrollContent(sim, "top");
		const cursor = sim.editor.locator(".cm-cursor-primary");

		await page.keyboard.press("ControlOrMeta+End");
		await page.keyboard.type("x");
		expect((await box(cursor)).bottom).toBeLessThanOrEqual(
			(await box(sim.runBar)).top + 1,
		);

		await page.keyboard.press("ControlOrMeta+Home");
		expect((await box(cursor)).top).toBeGreaterThanOrEqual(
			(await box(sim.toolbar)).bottom - 1,
		);
	});

	test("an invalid config shows its errors above the code and scrolls with it", async ({
		page,
	}) => {
		const sim = await openWith(
			page,
			{ width: 1440, height: 700 },
			`${longConfig}\nnot a valid line`,
		);
		await scrollContent(sim, "top");
		const panel = sim.editor.locator(".cm-panels-top [role=alert]");
		await expect(panel).toBeVisible({ timeout: 30_000 });
		const before = await box(panel);
		const code = await box(sim.editor.locator(".cm-content"));
		expect(before.bottom).toBeLessThanOrEqual(code.top + 1);
		expect(before.bottom).toBeLessThanOrEqual((await box(sim.runBar)).top);

		const scrolled = await sim.content.evaluate((el) => {
			el.scrollTop = 200;
			return el.scrollTop;
		});
		expect((await box(panel)).top).toBeCloseTo(before.top - scrolled, 0);
	});
});
