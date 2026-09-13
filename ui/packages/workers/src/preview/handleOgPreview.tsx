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

// Renders SatoriPreviewCard to PNG (or SVG for debugging) from a share/db key,
// as a manual preview aid. Distinct from the live OG path (handlePreview /
// /api/preview/*), which is left unchanged; cutting live OG over to this
// renderer is a separate follow-up.
export async function handleOgPreview(
	request: IRequest,
	env: Env,
): Promise<Response> {
	const { params } = request;
	const key = params?.key;
	if (!key) {
		return new Response("missing key", { status: 400 });
	}

	// Match handleView's db-variant detection: /api/og-preview/db/:key resolves
	// against the backend's /api/share/db/ path.
	const dbStr = request.url.includes("/db/") ? "db/" : "";

	// PNG is the default and only OG-valid output; ?format=svg is a browser-only
	// debug aid (unfurlers ignore SVG).
	const format =
		new URL(request.url).searchParams.get("format") === "svg" ? "svg" : "png";

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
		width: CARD_W,
		height: CARD_H,
		format,
		fonts: cardFonts,
	};
	return new ImageResponse(<SatoriPreviewCard data={result} />, options);
}
