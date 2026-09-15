import { expect, type Locator, type Page } from "@playwright/test";
import { ConsoleMonitor } from "./console-monitor";

/**
 * The reusable browser-driving harness for the gcsim **docs** site
 * (`@gcsim/docs`, Docusaurus).
 *
 * Docs is a single, uniform surface — every route is a rendered MDX page with
 * the same sidebar chrome — so, unlike the web app's per-route page objects,
 * one class owns the whole surface and bundles a {@link ConsoleMonitor}.
 * Docusaurus ships no data-testids, so every locator rides an observable
 * contract: the sidebar's `aria-label`, the theme's doc-body class, or an
 * accessible name.
 */
export class DocsHarness {
	readonly page: Page;
	readonly console: ConsoleMonitor;
	/** The docs sidebar (`<nav aria-label="Docs sidebar">`). */
	readonly sidebar: Locator;
	/** The rendered MDX body of the current doc (`article .theme-doc-markdown`). */
	readonly content: Locator;

	constructor(page: Page) {
		this.page = page;
		this.console = new ConsoleMonitor(page);
		this.sidebar = page.getByRole("navigation", { name: "Docs sidebar" });
		this.content = page.locator("article .theme-doc-markdown");
	}

	/** A sidebar entry by its exact visible label. */
	sidebarLink(name: string): Locator {
		return this.sidebar.getByRole("link", { name, exact: true });
	}

	/** The current doc's `<h1>`. */
	heading(name: string): Locator {
		return this.page.getByRole("heading", { level: 1, name });
	}

	/** Navigate to `path` and wait for the doc body to render. */
	async goto(path = "/"): Promise<void> {
		await this.page.goto(path);
		await expect(this.content).toBeVisible();
	}

	/**
	 * Docusaurus marks doc images `loading="lazy"`, so off-screen ones never
	 * fetch on their own — each is scrolled into view first to start the load.
	 */
	async assertImagesLoaded(container: Locator): Promise<void> {
		const images = container.locator("img");
		const count = await images.count();
		for (let i = 0; i < count; i++) {
			await images.nth(i).scrollIntoViewIfNeeded();
		}
		await expect(async () => {
			const broken = await images.evaluateAll((els) =>
				(els as HTMLImageElement[])
					.filter((img) => !img.complete || img.naturalWidth === 0)
					.map((img) => img.currentSrc || img.src),
			);
			expect(broken, `images failed to load:\n${broken.join("\n")}`).toEqual(
				[],
			);
		}).toPass();
	}
}
