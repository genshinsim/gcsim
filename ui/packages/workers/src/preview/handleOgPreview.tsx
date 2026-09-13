import type { SatoriFont } from "@gcsim/components/src/Cards/SatoriPreviewCard/fonts";
// Deep import (not the package root) so bundling pulls only the card subtree,
// not the whole component library (which imports .png/.css assets).
import { SatoriPreviewCard } from "@gcsim/components/src/Cards/SatoriPreviewCard/SatoriPreviewCard";
import type { model } from "@gcsim/types";
import type { IRequest } from "itty-router";
import { ImageResponse } from "workers-og";
import type { Env } from "../bindings";
import { cardFonts } from "./cardFonts";

// workers-og 0.0.27 derives its option type from @vercel/og, which it does not
// depend on; the resulting `ImageResponseOptions` drops the `fonts` field satori
// supports at runtime. Reconstruct a usable shape from the constructor and add
// `fonts` back so the render is fully typed.
type ImageResponseOptions = NonNullable<
	ConstructorParameters<typeof ImageResponse>[1]
> & { fonts?: SatoriFont[] };

// Card geometry, matching SatoriPreviewCard's fixed layout and the live OG
// dimensions (og:image:width/height in handleInjectHead).
const CARD_W = 540;
const CARD_H = 250;

// Supersample factor. workers-og rasterises satori's SVG at 1:1 with the card's
// logical width, so at the card's native 540x250 the PNG is soft once an
// unfurler scales it up. Rendering the (vector) card at 2x — 1080x500 — and
// letting clients display it at the 540x250 the og:image meta advertises keeps
// the raster crisp. Pure resolution bump: the card layout is unchanged, only
// wrapped in a scaled container.
const SCALE = 2;

// Long-lived edge cache. Satori rendering is far more expensive than the old
// backend PNG proxy, so render-once/serve-many matters more here, not less.
// Matches the retired handlePreview proxy's TTL (60 days).
const CACHE_TTL_SECONDS = 60 * 24 * 60 * 60;

// Renders SatoriPreviewCard to PNG (or SVG for debugging) from a share/db key.
// This is the live OG path: `/api/preview/:key` and `/api/preview/db/:key`, hit
// by crawlers via the og:image meta (`.../api/preview/<key>.png`). Renders once
// and serves the result from caches.default thereafter.
export async function handleOgPreview(
	request: IRequest,
	env: Env,
	ctx: ExecutionContext,
): Promise<Response> {
	const { params } = request;
	// The og:image meta points crawlers at `/api/preview/<key>.png`, but the key
	// is also valid without the suffix — strip `.png` either way so the backend
	// share lookup (which has no such suffix) resolves.
	const key = params?.key?.replace(/\.png$/, "");
	if (!key) {
		return new Response("missing key", { status: 400 });
	}

	const url = new URL(request.url);

	// Match handleView's db-variant detection: /api/preview/db/:key resolves
	// against the backend's /api/share/db/ path. Test the pathname (not the whole
	// URL) so a `?...=/db/` query string can't flip the variant.
	const dbStr = url.pathname.includes("/db/") ? "db/" : "";

	// PNG is the default and only OG-valid output; ?format=svg is a browser-only
	// debug aid (unfurlers ignore SVG).
	const format = url.searchParams.get("format") === "svg" ? "svg" : "png";

	// Serve a previously rendered card from the edge cache when present. Keyed on
	// the full request URL so the `.png`/no-suffix and format variants stay
	// distinct.
	const cache = caches.default;
	const cacheKey = new Request(request.url, request);
	const cached = await cache.match(cacheKey);
	if (cached) {
		return cached;
	}

	// Same data source as the live preview/share path: the backend share JSON.
	const resp = await fetch(
		new Request(env.API_ENDPOINT + "/api/share/" + dbStr + key),
	);
	if (!resp.ok) {
		// Unknown/invalid key must be a 4xx, not a 500. Pass through the backend's
		// 4xx (e.g. 404 not found); collapse any 5xx to 502 Bad Gateway.
		const status = resp.status >= 400 && resp.status < 500 ? resp.status : 502;
		return new Response("could not load share " + key, { status });
	}

	let result: model.SimulationResult;
	try {
		result = (await resp.json()) as model.SimulationResult;
	} catch {
		return new Response("invalid share data for " + key, { status: 400 });
	}

	const options: ImageResponseOptions = {
		width: CARD_W * SCALE,
		height: CARD_H * SCALE,
		format,
		fonts: cardFonts,
	};
	const image = new ImageResponse(
		// Scale the fixed-size card up to fill the supersampled canvas without
		// touching its own layout.
		<div
			style={{
				display: "flex",
				transform: `scale(${SCALE})`,
				transformOrigin: "top left",
			}}
		>
			<SatoriPreviewCard data={result} />
		</div>,
		options,
	);

	// Cache the rendered card (workers-og already sets a long immutable
	// Cache-Control; caches.default is what actually spares the re-render). Only
	// PNGs are cached — SVG is a debug-only path.
	if (format === "png") {
		const cacheable = new Response(image.body, image);
		cacheable.headers.set("Cache-Control", `max-age=${CACHE_TTL_SECONDS}`);
		ctx.waitUntil(cache.put(cacheKey, cacheable.clone()));
		return cacheable;
	}

	return image;
}
