import type { Env } from "../bindings";

// A 1x1 transparent PNG. Used in place of any asset the Worker could not fetch
// so the render stays valid (Satori rejects a non-absolute <img src>) and makes
// no network fetch of its own. A render that had to use it is never cached.
const FALLBACK_ASSET =
	"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==";

// Base64-encode arbitrary bytes for a `data:` URI. Chunked so large buffers
// don't blow the argument limit of String.fromCharCode; btoa is available in
// the Workers runtime.
function base64FromArrayBuffer(buffer: ArrayBuffer): string {
	const bytes = new Uint8Array(buffer);
	let binary = "";
	const chunk = 0x8000;
	for (let i = 0; i < bytes.length; i += chunk) {
		binary += String.fromCharCode(...bytes.subarray(i, i + chunk));
	}
	return btoa(binary);
}

export type ResolvedAssets = {
	// Maps a relative asset path (e.g. "avatar/Nahida.png") to an <img src>: the
	// fetched bytes as a `data:` URI, or the placeholder for a path that failed.
	resolve: (path: string) => string;
	// True when any asset failed to fetch and the fallback stood in — the render
	// must not be cached so it re-renders once the asset lands.
	usedFallback: boolean;
};

// Pre-fetch the card's assets in the Worker and return a resolver that hands
// back `data:` URIs, so Satori makes zero outbound image fetches at render
// time. Paths are the relative ones enumerateAssetPaths yields; each is fetched
// once (callers dedup) from the same origin the /api/assets proxy uses.
export async function fetchCardAssets(
	paths: string[],
	env: Env,
): Promise<ResolvedAssets> {
	const map = new Map<string, string>();
	let usedFallback = false;

	await Promise.all(
		paths.map(async (path) => {
			try {
				const resp = await fetch(
					new Request(`${env.ASSETS_ENDPOINT}/api/assets/${path}`),
				);
				if (!resp.ok) {
					usedFallback = true;
					return;
				}
				const buffer = await resp.arrayBuffer();
				const contentType = resp.headers.get("Content-Type") ?? "image/png";
				map.set(
					path,
					`data:${contentType};base64,${base64FromArrayBuffer(buffer)}`,
				);
			} catch {
				usedFallback = true;
			}
		}),
	);

	return {
		resolve: (path) => map.get(path) ?? FALLBACK_ASSET,
		usedFallback,
	};
}
