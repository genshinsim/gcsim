import "@testing-library/jest-dom/vitest";

// jsdom does not implement ResizeObserver, which cmdk (and Radix components
// built on it) use to track content size; stub it so their effects can run.
class ResizeObserverStub {
	observe() {}
	unobserve() {}
	disconnect() {}
}

if (typeof globalThis.ResizeObserver === "undefined") {
	globalThis.ResizeObserver =
		ResizeObserverStub as unknown as typeof ResizeObserver;
}

// jsdom also doesn't implement scrollIntoView, which cmdk calls when the
// keyboard-highlighted item changes.
if (typeof Element.prototype.scrollIntoView === "undefined") {
	Element.prototype.scrollIntoView = () => {};
}
