import type { Env } from "../bindings";
import { newShareKey } from "./shareKey";

const MAX_ATTEMPTS = 5;

async function gzip(data: ArrayBuffer): Promise<ArrayBuffer> {
	const stream = new Blob([data])
		.stream()
		.pipeThrough(new CompressionStream("gzip"));
	return new Response(stream).arrayBuffer();
}

// Stores data gzipped under a new key built from prefix and returns the key.
export async function storeShare(
	env: Env,
	prefix: string,
	data: ArrayBuffer,
): Promise<string> {
	const gz = await gzip(data);
	for (let i = 0; i < MAX_ATTEMPTS; i++) {
		const key = newShareKey(prefix);
		const obj = await env.GCSIM_SHARES.put(key, gz, {
			onlyIf: new Headers({ "If-None-Match": "*" }),
			httpMetadata: {
				contentType: "application/json",
				contentEncoding: "gzip",
			},
		});
		if (obj != null) {
			return key;
		}
	}
	throw new Error("no free share key");
}

// Looks a share up in R2, falling back to the backend on a miss.
async function findShare(
	env: Env,
	key: string,
): Promise<R2ObjectBody | Response> {
	const obj = await env.GCSIM_SHARES.get(key);
	return obj ?? fetch(new Request(`${env.API_ENDPOINT}/api/share/${key}`));
}

// Response for a client: a stored share is sent as its gzip bytes, a backend share is gzipped on the way out.
export async function shareResponse(env: Env, key: string): Promise<Response> {
	const found = await findShare(env, key);
	if (found instanceof Response) {
		const res = new Response(found.body, found);
		res.headers.append("Content-Encoding", "gzip");
		return res;
	}
	return new Response(found.body, {
		headers: {
			"Content-Type": "application/json",
			"Content-Encoding": "gzip",
		},
		encodeBody: "manual",
	});
}

// Decoded share JSON as a Response, for reading inside the worker.
export async function shareJSON(env: Env, key: string): Promise<Response> {
	const found = await findShare(env, key);
	if (found instanceof Response) {
		return found;
	}
	return new Response(found.body.pipeThrough(new DecompressionStream("gzip")));
}
