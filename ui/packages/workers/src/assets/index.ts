import type { IRequest } from "itty-router";
import type { Env } from "../bindings";
import artifactMap from "./artifact.dm.json";
import avatarMap from "./character.dm.json";
import weaponMap from "./weapon.dm.json";

// Asset request handled entirely by the worker (ported from the Go assets
// microservice in internal/services/assets). Layout of the GCSIM_ASSETS bucket
// mirrors that service's routes: dynamic images are cached under
// `<type>/<key>.png`, and the static files (fonts, logo, special/*, misc/*)
// live at their sub-path, uploaded on deploy.

type AssetType = "avatar" | "weapons" | "artifacts";

const NAME_MAPS: Record<AssetType, Record<string, string>> = {
	avatar: avatarMap,
	weapons: weaponMap,
	artifacts: artifactMap,
};

// Fallback source hosts if ASSET_SOURCE_HOSTS is unset/invalid. Kept in sync
// with the Go service's mihoyo defaults (cmd/services/assets/main.go).
const DEFAULT_SOURCE_HOSTS: Record<AssetType, string[]> = {
	avatar: [
		"https://upload-os-bbs.mihoyo.com/game_record/genshin/character_icon/",
	],
	weapons: ["https://upload-os-bbs.mihoyo.com/game_record/genshin/equip/"],
	artifacts: ["https://upload-os-bbs.mihoyo.com/game_record/genshin/equip/"],
};

// Traveler variants ship as bundled static icons (static/special/<key>.png)
// rather than a mihoyo CDN filename. Mirrors internal/services/assets/special.go.
const SPECIAL_KEYS = new Set<string>([
	"aetheranemo",
	"lumineanemo",
	"aetherdendro",
	"luminedendro",
	"aetherelectro",
	"lumineelectro",
	"aethergeo",
	"luminegeo",
	"aetherhydro",
	"luminehydro",
]);

const DEFAULT_KEY = "misc/default.png";
const CACHE_SECONDS = 60 * 24 * 60 * 60; // 60 days
const LONG_CACHE = `max-age=${CACHE_SECONDS}`;
const TYPED_RE = /^(avatar|weapons|artifacts)\/(.+)\.png$/;

export async function handleAssets(
	request: IRequest,
	env: Env,
	ctx: ExecutionContext,
): Promise<Response> {
	const cacheUrl = new URL(request.url);
	const cacheKey = new Request(cacheUrl.toString(), request);
	const cache = caches.default;
	const cached = await cache.match(cacheKey);
	if (cached) {
		return cached;
	}

	const subpath = cacheUrl.pathname.replace(/^\/api\/assets\//, "");
	const typed = subpath.match(TYPED_RE);

	const response = typed
		? await resolveTyped(typed[1] as AssetType, typed[2], subpath, env, ctx)
		: await serveStatic(subpath, env);

	// Never persist the no-cache fallback for a requested key, so a later
	// successful fetch can still populate it (matches the Go service).
	if (
		response.status === 200 &&
		response.headers.get("Cache-Control") !== "no-cache"
	) {
		ctx.waitUntil(cache.put(cacheKey, response.clone()));
	}
	return response;
}

async function resolveTyped(
	type: AssetType,
	key: string,
	subpath: string,
	env: Env,
	ctx: ExecutionContext,
): Promise<Response> {
	// Travelers: serve the bundled special icon, not the mihoyo CDN image.
	// Checked for every type before the map lookup, matching the Go service.
	if (SPECIAL_KEYS.has(key)) {
		const special = await env.GCSIM_ASSETS.get(`special/${key}.png`);
		return special ? r2Response(special, `special/${key}.png`) : fallback(env);
	}

	// R2 cache hit.
	const cached = await env.GCSIM_ASSETS.get(subpath);
	if (cached) {
		return r2Response(cached, subpath);
	}

	// Resolve the mihoyo CDN filename; unknown key -> default fallback.
	const assetName = NAME_MAPS[type][key];
	if (!assetName) {
		return fallback(env);
	}

	// Try each source host in order; first valid image wins and is cached to R2.
	for (const host of sourceHosts(type, env)) {
		const image = await fetchImage(host, assetName);
		if (!image) {
			continue;
		}
		ctx.waitUntil(
			env.GCSIM_ASSETS.put(subpath, image.body, {
				httpMetadata: {
					contentType: image.contentType,
					cacheControl: LONG_CACHE,
				},
			}),
		);
		return new Response(image.body, {
			headers: {
				"Content-Type": image.contentType,
				"Cache-Control": LONG_CACHE,
			},
		});
	}

	return fallback(env);
}

function sourceHosts(type: AssetType, env: Env): string[] {
	if (env.ASSET_SOURCE_HOSTS) {
		try {
			const parsed = JSON.parse(env.ASSET_SOURCE_HOSTS) as Partial<
				Record<AssetType, string[]>
			>;
			const hosts = parsed[type];
			if (hosts && hosts.length > 0) {
				return hosts;
			}
		} catch {
			// fall through to hardcoded defaults
		}
	}
	return DEFAULT_SOURCE_HOSTS[type];
}

async function fetchImage(
	host: string,
	assetName: string,
): Promise<{ body: ArrayBuffer; contentType: string } | null> {
	const url = new URL(`${assetName}.png`, ensureTrailingSlash(host)).toString();
	let resp: Response;
	try {
		resp = await fetch(url, {
			cf: { cacheTtl: CACHE_SECONDS, cacheEverything: true },
		});
	} catch {
		return null;
	}
	if (resp.status !== 200) {
		return null;
	}
	const contentType = resp.headers.get("Content-Type") ?? "";
	if (!contentType.startsWith("image/")) {
		return null;
	}
	return { body: await resp.arrayBuffer(), contentType };
}

// Serve a static file straight from R2; a miss is a 404 (matches the Go file
// server mounted at /api/assets/*).
async function serveStatic(key: string, env: Env): Promise<Response> {
	const object = await env.GCSIM_ASSETS.get(key);
	if (!object) {
		return new Response("Not Found", { status: 404 });
	}
	return r2Response(object, key);
}

function r2Response(object: R2ObjectBody, key: string): Response {
	const headers = new Headers();
	object.writeHttpMetadata(headers);
	headers.set("etag", object.httpEtag);
	if (!headers.has("Cache-Control")) {
		headers.set("Cache-Control", LONG_CACHE);
	}
	if (!headers.has("Content-Type")) {
		headers.set("Content-Type", contentTypeFor(key));
	}
	return new Response(object.body, { headers });
}

async function fallback(env: Env): Promise<Response> {
	const object = await env.GCSIM_ASSETS.get(DEFAULT_KEY);
	if (!object) {
		return new Response("Not Found", { status: 404 });
	}
	return new Response(object.body, {
		headers: {
			"Content-Type": "image/png",
			"Cache-Control": "no-cache",
		},
	});
}

function ensureTrailingSlash(host: string): string {
	return host.endsWith("/") ? host : `${host}/`;
}

function contentTypeFor(key: string): string {
	if (key.endsWith(".png")) return "image/png";
	if (key.endsWith(".jpg") || key.endsWith(".jpeg")) return "image/jpeg";
	if (key.endsWith(".ttf")) return "font/ttf";
	if (key.endsWith(".svg")) return "image/svg+xml";
	return "application/octet-stream";
}
