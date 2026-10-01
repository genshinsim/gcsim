// CodeMirror measures text through Range rects, which jsdom doesn't
// implement. Import this from the test setup of any package whose tests
// render the editor under jsdom.
if (typeof Range !== "undefined") {
	Range.prototype.getClientRects ??= () =>
		Object.assign([], { item: () => null }) as unknown as DOMRectList;
	Range.prototype.getBoundingClientRect ??= () => new DOMRect();
}
