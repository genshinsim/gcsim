import type { IRequest } from "itty-router";
import type { Env } from "../bindings";

export async function handleDBRedirect(
	request: IRequest,
	env: Env,
): Promise<Response> {
	const shareKey = await lookupShareKey(env, request.params.key);
	if (!shareKey) {
		return env.ASSETS.fetch(request);
	}
	const url = new URL(request.url);
	url.pathname = `/sh/${shareKey}`;
	return Response.redirect(url.toString(), 301);
}

async function lookupShareKey(
	env: Env,
	id: string,
): Promise<string | undefined> {
	try {
		const res = await fetch(
			`${env.API_ENDPOINT}/api/dbshare/${encodeURIComponent(id)}`,
		);
		if (!res.ok) return undefined;
		const { share_key } = (await res.json()) as { share_key?: string };
		return share_key || undefined;
	} catch (err) {
		console.log(`db share key lookup failed for ${id}: ${err}`);
		return undefined;
	}
}
