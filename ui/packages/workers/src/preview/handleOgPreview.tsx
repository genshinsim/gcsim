import { enumerateAssetPaths } from "@gcsim/components/src/Cards/SatoriPreviewCard/assetPaths";
import type { SatoriFont } from "@gcsim/components/src/Cards/SatoriPreviewCard/fonts";
import { compositePortraits } from "@gcsim/components/src/Cards/SatoriPreviewCard/portraitCompositor";
// Deep import (not the package root) so bundling pulls only the card subtree,
// not the whole component library (which imports .png/.css assets).
import { SatoriPreviewCard } from "@gcsim/components/src/Cards/SatoriPreviewCard/SatoriPreviewCard";
import type { model } from "@gcsim/types";
import type { IRequest } from "itty-router";
import { ImageResponse } from "workers-og";
import type { Env } from "../bindings";
import { cardFonts } from "./cardFonts";
import { fetchCardAssets } from "./fetchCardAssets";

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
// unfurler scales it up. Rendering the (vector) card above native — e.g. 2x,
// 1080x500 — and letting clients display it at the 540x250 the og:image meta
// advertises keeps the raster crisp and gives Discord (which downscales the
// source) a sharper image. Portraits are composited at this same scale so they
// aren't the soft bottleneck. Source assets are 256x256 and avatars logical
// 96px, so higher factors mildly interpolate avatars while weapons/flowers
// never upscale — an accepted tradeoff for crisper vector text/chart edges.
const DEFAULT_SCALE = 1.5;
const MIN_SCALE = 1; // below native only softens the image
const MAX_SCALE = 3; // cap render latency / canvas size

// Resolve the supersample factor from OG_PREVIEW_SCALE (dashboard-managed, like
// ASSET_SOURCE_HOSTS). Unset, non-numeric, or garbage falls back to the default
// without throwing; the result is clamped to [MIN_SCALE, MAX_SCALE].
function resolveScale(env: Env): number {
	const parsed = parseFloat(env.OG_PREVIEW_SCALE ?? "");
	const scale = Number.isNaN(parsed) ? DEFAULT_SCALE : parsed;
	return Math.min(MAX_SCALE, Math.max(MIN_SCALE, scale));
}

// Edge-cache revision. Mixed into the cache key (see below) so bumping it makes
// every cached card miss and re-render — the manual invalidation lever for
// whenever rendering or the scale is tweaked (e.g. after changing
// OG_PREVIEW_SCALE, since the scale itself is deliberately not in the key).
// Read from OG_PREVIEW_CACHE_REV (dashboard-managed, like OG_PREVIEW_SCALE) so
// it can be bumped without a code change; unset/invalid falls back to the
// hardcoded default. Plain incrementing integer, not semver.
const DEFAULT_CACHE_REV = 1;

function resolveCacheRev(env: Env): number {
	const parsed = parseInt(env.OG_PREVIEW_CACHE_REV ?? "", 10);
	return Number.isNaN(parsed) ? DEFAULT_CACHE_REV : parsed;
}

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
	// the full request URL (so the `.png`/no-suffix and format variants stay
	// distinct) plus the cache-rev knob via a synthetic `_ogrev` query param, so
	// bumping the rev partitions the cache into a fresh namespace. The resolved
	// scale is deliberately not part of the key — the rev is the invalidation
	// lever.
	const cache = caches.default;
	const cacheUrl = new URL(request.url);
	cacheUrl.searchParams.set("_ogrev", String(resolveCacheRev(env)));
	const cacheKey = new Request(cacheUrl.toString(), request);
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

	// Pre-fetch the card's imagery here (deduped, in parallel) through the
	// Worker's own asset resolver, so the compositor (and thus Satori) makes zero
	// outbound image fetches. A missing asset resolves to the misc/default.png
	// placeholder and marks the render uncacheable, so a card with a missing asset
	// isn't frozen for the full TTL.
	const assets = await fetchCardAssets(enumerateAssetPaths(result), env, ctx);

	// Pre-composite each portrait (bg + avatar + weapon + artifact set(s)) into
	// one flat PNG per slot, with the white icon outline and two-set slice baked
	// in at raster level — Satori can express neither. Runs at the supersample
	// scale so portraits match the card's crispness.
	const scale = resolveScale(env);
	const portraits = compositePortraits(result, assets.resolveBytes, scale);

	const options: ImageResponseOptions = {
		width: CARD_W * scale,
		height: CARD_H * scale,
		format,
		fonts: cardFonts,
	};
	const image = new ImageResponse(
		// Scale the fixed-size card up to fill the supersampled canvas without
		// touching its own layout.
		<div
			style={{
				display: "flex",
				transform: `scale(${scale})`,
				transformOrigin: "top left",
			}}
		>
			<SatoriPreviewCard data={result} portraits={portraits} />
		</div>,
		options,
	);

	// A render that fell back to a placeholder must not be stored anywhere: skip
	// caches.default and mark it no-store so the next request re-renders once the
	// missing asset lands. It is still served (200 image/png).
	if (assets.usedFallback) {
		const uncacheable = new Response(image.body, image);
		uncacheable.headers.set("Cache-Control", "no-store");
		return uncacheable;
	}

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
