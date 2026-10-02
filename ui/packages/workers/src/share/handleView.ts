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

	//check if this is db route
	const isDB = request.url.includes("/db/");

	let response = await cache.match(cacheKey);

	if (!response) {
		console.log(
			`Response for request url: ${request.url} not present in cache. Fetching and caching request.`,
		);

		if (isDB) {
			response = await fetch(
				new Request(`${env.API_ENDPOINT}/api/share/db/${key}`),
			);
			response = new Response(response.body, response);
			response.headers.append("Content-Encoding", "gzip");
		} else {
			response = await shareResponse(env, key);
		}
		response.headers.append("Cache-Control", "s-maxage=1800");

		ctx.waitUntil(cache.put(cacheKey, response.clone()));
	} else {
		console.log(`cache hit for: ${request.url}`);
	}

	return response;
}
