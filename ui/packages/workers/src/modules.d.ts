// TTF files are bundled as raw bytes via the wrangler `Data` module rule
// (see wrangler.jsonc); the default export is the file's ArrayBuffer.
declare module "*.ttf" {
	const data: ArrayBuffer;
	export default data;
}
