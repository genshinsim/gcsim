import { resolveAsset, STATUS_HEADER } from "../assets";
import type { Env } from "../bindings";

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
	// resolved bytes as a `data:` URI. A missing asset resolves to the Worker's
	// misc/default.png placeholder bytes (same as /api/assets/* serves).
	resolve: (path: string) => string;
	// True when any asset resolved to the placeholder (X-Gcsim-Asset: fallback) —
	// the render must not be cached so it re-renders once the asset lands.
	usedFallback: boolean;
};

// Pre-fetch the card's assets in the Worker and return a resolver that hands
// back `data:` URIs, so Satori makes zero outbound image fetches at render
// time. Each path is resolved once (callers dedup) through the Worker's own
// in-process asset resolver — the same resolution the public /api/assets/*
// route uses — so there is no outbound HTTP round-trip and no duplicated
// R2/source-host/placeholder logic.
export async function fetchCardAssets(
	paths: string[],
	env: Env,
	ctx: ExecutionContext,
): Promise<ResolvedAssets> {
	const map = new Map<string, string>();
	let usedFallback = false;

	await Promise.all(
		paths.map(async (path) => {
			let resp: Response;
			try {
				resp = await resolveAsset(path, env, ctx);
			} catch {
				usedFallback = true;
				return;
			}
			if (resp.status !== 200) {
				usedFallback = true;
				return;
			}
			// A genuinely-missing asset comes back as the misc/default.png
			// placeholder (HTTP 200), flagged X-Gcsim-Asset: fallback — not by a
			// non-ok status. That is the signal to skip caching this render.
			if (resp.headers.get(STATUS_HEADER) === "fallback") {
				usedFallback = true;
			}
			const buffer = await resp.arrayBuffer();
			const contentType = resp.headers.get("Content-Type") ?? "image/png";
			map.set(
				path,
				`data:${contentType};base64,${base64FromArrayBuffer(buffer)}`,
			);
		}),
	);

	// Every enumerated path maps to real bytes here: a typed miss resolves to the
	// placeholder (200), and the static paths are bundled files. An empty src can
	// only occur if a bundled static file (nahida/default) is itself absent — a
	// broken deploy — and such a render is uncacheable (usedFallback) regardless.
	return {
		resolve: (path) => map.get(path) ?? "",
		usedFallback,
	};
}
