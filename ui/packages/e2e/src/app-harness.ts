import type { Page } from "@playwright/test";
import { ConsoleMonitor } from "./console-monitor";
import { SimulatorPage } from "./pages/simulator-page";
import { ViewerPage } from "./pages/viewer-page";

/**
 * The reusable browser-driving harness for the gcsim web app.
 *
 * It encodes the app's implicit protocol once — boot, wait for wasm+workers,
 * load a config, run a sim, open the viewer — as a standalone importable layer.
 * The regression specs consume it; a future agent-exploration head can too.
 */
export class AppHarness {
	readonly page: Page;
	readonly simulator: SimulatorPage;
	readonly viewer: ViewerPage;
	readonly console: ConsoleMonitor;

	constructor(page: Page) {
		this.page = page;
		this.simulator = new SimulatorPage(page);
		this.viewer = new ViewerPage(page);
		this.console = new ConsoleMonitor(page);
	}

	/** Open the simulator and wait for wasm + workers to be ready. */
	async boot(): Promise<void> {
		await this.simulator.goto();
		await this.simulator.waitForReady();
	}

	/**
	 * Run `cfg` and wait for the run to complete.
	 *
	 * Completion is the `run time: N ms` console signal. It is matched against
	 * the {@link ConsoleMonitor} buffer rather than a fresh `waitForEvent`,
	 * because a one-iteration sim can log it before a post-click listener would
	 * attach.
	 */
	async run(cfg: string): Promise<void> {
		await this.simulator.setConfig(cfg);
		await this.simulator.waitForConfigValid();
		await this.simulator.run();
		// Anchored so a cancel ("cancelled with run time: N ms") can't satisfy it.
		await this.console.waitForLog(/^run time: \d+(\.\d+)? ms$/);
	}
}
