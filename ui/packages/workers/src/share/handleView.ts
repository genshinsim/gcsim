import type { IRequest } from "itty-router";
import type { Env } from "../bindings";
import { shareResponse } from "./storage";

export async function handleView(
	request: IRequest,
	env: Env,
	ctx: ExecutionContext,
): Promise<Response> {
	const { params } = request;
	if (!params || !params.key) {
		return new Response(null, {
			status: 400,
			statusText: "Bad Request",
		});
	}

	const key = params.key;

	if (key === "") {
		return new Response(null, {
			status: 400,
			statusText: "Bad Request",
		});
	}

	console.log(key);

	const cacheUrl = new URL(request.url);
	const cacheKey = new Request(cacheUrl.toString(), request);
	console.log(`checking for cache key: ${cacheUrl}`);
	const cache = caches.default;

	let response = await cache.match(cacheKey);

	if (!response) {
		console.log(
			`Response for request url: ${request.url} not present in cache. Fetching and caching request.`,
		);

		response = await shareResponse(env, key);
		if (response.status === 200) {
			response.headers.set("Cache-Control", "max-age=14400, s-maxage=31536000");
			ctx.waitUntil(cache.put(cacheKey, response.clone()));
		}
	} else {
		console.log(`cache hit for: ${request.url}`);
	}

	return response;
}
