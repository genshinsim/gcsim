// Minimal ambient type for the one Node builtin fontsNode.ts uses, so the
// browser-oriented @gcsim/components package needn't depend on @types/node
// (which would leak Node globals into the whole package's typecheck). This file
// has no imports/exports, so `declare module` declares the ambient module
// rather than augmenting a non-existent one.
declare module "node:fs/promises" {
	export function readFile(path: URL): Promise<Uint8Array>;
}
