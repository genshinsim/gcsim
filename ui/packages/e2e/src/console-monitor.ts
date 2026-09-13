import { expect, type Page } from "@playwright/test";

// Messages that are environmental noise rather than real app faults. Keep this
// list short and justified — every entry weakens the "zero console errors"
// guarantee. Match against the full console/error text.
const IGNORED: RegExp[] = [
	// Vite dev server injects a websocket HMR client; a transient connection
	// blip during teardown is not an app fault.
	/\[vite\] failed to connect to websocket/i,
	// React routes its dev-only warnings through console.error, prefixed
	// "Warning:". Under StrictMode the app + its third-party libs (react-helmet,
	// react-transition-group, blueprint) emit a fixed set — legacy lifecycles,
	// findDOMNode, the language <select>'s `selected` option. These are
	// pre-existing and dev-only (absent from a production build); a real crash
	// surfaces as a `pageerror` or a non-"Warning:" console.error, which this
	// does NOT mask.
	/^Warning: /,
	// visx draws a couple of markers with a computed negative radius on this
	// fixture (e.g. a zero-variance distribution), so the browser rejects
	// `<circle r="-2">`. Pre-existing chart quirk, unrelated to the boot/run/
	// render path this smoke test guards.
	/<circle> attribute r: A negative value/,
];

function isIgnored(text: string): boolean {
	return IGNORED.some((re) => re.test(text));
}

type Waiter = (text: string) => void;

/**
 * Watches a page's console for the lifetime of the page.
 *
 * Two jobs:
 *  - collect uncaught `console.error` / `pageerror` so the smoke spec can assert
 *    the critical path ran clean;
 *  - buffer every console message so {@link waitForLog} can match a signal even
 *    if it already fired. A one-iteration sim finishes in well under a second,
 *    so a post-hoc `page.waitForEvent` would race the message it waits for.
 */
export class ConsoleMonitor {
	readonly errors: string[] = [];
	private readonly logs: string[] = [];
	private readonly waiters = new Set<Waiter>();

	constructor(page: Page) {
		page.on("console", (msg) => {
			const text = msg.text();
			this.record(text);
			if (msg.type() === "error" && !isIgnored(text)) {
				this.errors.push(`console.error: ${text}`);
			}
		});
		page.on("pageerror", (err) => {
			const text = err.stack ?? err.message;
			this.record(text);
			if (!isIgnored(text)) {
				this.errors.push(`pageerror: ${text}`);
			}
		});
	}

	private record(text: string): void {
		this.logs.push(text);
		for (const waiter of this.waiters) {
			waiter(text);
		}
	}

	/**
	 * Resolve once a console message matching `re` has been seen — checking the
	 * buffer first, so an already-fired message counts. Rejects on timeout with
	 * the recent buffer attached for diagnosis.
	 */
	waitForLog(re: RegExp, timeout = 60_000): Promise<string> {
		const existing = this.logs.find((text) => re.test(text));
		if (existing !== undefined) {
			return Promise.resolve(existing);
		}
		return new Promise((resolve, reject) => {
			const waiter: Waiter = (text) => {
				if (re.test(text)) {
					cleanup();
					resolve(text);
				}
			};
			const timer = setTimeout(() => {
				cleanup();
				const tail = this.logs.slice(-20).join("\n");
				reject(
					new Error(
						`timed out after ${timeout}ms waiting for console log matching ${re}.\nrecent console:\n${tail}`,
					),
				);
			}, timeout);
			const cleanup = () => {
				clearTimeout(timer);
				this.waiters.delete(waiter);
			};
			this.waiters.add(waiter);
		});
	}

	/** Fails the test if any un-ignored console/page error was seen so far. */
	assertNoErrors(): void {
		expect(
			this.errors,
			`expected no uncaught console errors, saw:\n${this.errors.join("\n")}`,
		).toEqual([]);
	}
}
