import type { IRequest } from "itty-router";
import type { Env } from "../bindings";
import { verifyShare } from "./auth";
import { storeShare } from "./storage";
import { validator } from "./validation";

const MAX_SHARE_BYTES = 1024 * 1024;
const AUTH_HEADER = "X-GCSIM-SHARE-AUTH";

type ShareBody = {
	key_type?: unknown;
};

function reject(status: number, statusText: string): Response {
	return new Response(null, { status, statusText });
}

export async function handleShare(
	request: IRequest,
	env: Env,
): Promise<Response> {
	if (Number(request.headers.get("Content-Length")) > MAX_SHARE_BYTES) {
		return reject(413, "Payload Too Large");
	}
	const body = await request.arrayBuffer();
	if (body.byteLength > MAX_SHARE_BYTES) {
		return reject(413, "Payload Too Large");
	}

	const auth = request.headers.get(AUTH_HEADER);
	if (!auth) {
		return reject(403, "Forbidden");
	}
	const [id, sig] = splitAuth(auth);
	if (id === "prod" || id === "dev") {
		return fetch(`${env.API_ENDPOINT}/api/share`, {
			method: "POST",
			body,
			headers: request.headers,
		});
	}

	const keyClass = await verifyShare(env.SHARE_KEYS, id, sig, body);
	if (keyClass == null) {
		return reject(400, "Bad Request");
	}

	let content: ShareBody;
	try {
		content = JSON.parse(new TextDecoder().decode(body));
	} catch {
		return reject(400, "Bad Request (Invalid JSON)");
	}
	if (!validator.validate(content).valid || content.key_type !== keyClass) {
		return reject(400, "Bad Request");
	}

	const key = await storeShare(env, body);
	return new Response(key, { status: 202 });
}

function splitAuth(auth: string): [string, string] {
	const i = auth.indexOf(":");
	return i === -1 ? [auth, ""] : [auth.slice(0, i), auth.slice(i + 1)];
}
