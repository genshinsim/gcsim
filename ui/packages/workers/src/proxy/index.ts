import type { IRequest } from "itty-router";
import type { Env } from "../bindings";

export async function proxyRequest(
	request: IRequest,
	env: Env,
): Promise<Response> {
	const x = new URL(request.url);
	return fetch(new Request(env.API_ENDPOINT + x.pathname + x.search, request));
}
