import "@testing-library/jest-dom/vitest";
import "@gcsim/editor/jsdom-setup";

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

// floating-ui (Radix popovers, selects, tooltips) checks for the top layer
// with matches(":popover-open") and matches(":modal") on every position
// update. jsdom's selector engine throws on both, and the throw-and-catch
// makes opening a popover take seconds; nothing is in the top layer here.
const nativeMatches = Element.prototype.matches;
Element.prototype.matches = function (this: Element, selector: string) {
	if (selector === ":popover-open" || selector === ":modal") return false;
	return nativeMatches.call(this, selector);
} as typeof nativeMatches;
