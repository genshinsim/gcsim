import type { Page } from "@playwright/test";
import { ConsoleMonitor } from "./console-monitor";
import { DbDatabasePage } from "./pages/db-database-page";
import { DbHomePage } from "./pages/db-home-page";

/**
 * The reusable browser-driving harness for the gcsim db ("Simpact") app.
 *
 * Mirrors {@link AppHarness}: it bundles the db page objects and the console
 * monitor. Network stubbing lives in {@link installDbRoutes}, wired in by the
 * `db` test fixture before the harness is handed to a spec.
 */
export class DbHarness {
	readonly page: Page;
	readonly home: DbHomePage;
	readonly database: DbDatabasePage;
	readonly console: ConsoleMonitor;

	constructor(page: Page) {
		this.page = page;
		this.home = new DbHomePage(page);
		this.database = new DbDatabasePage(page);
		this.console = new ConsoleMonitor(page);
	}
}
