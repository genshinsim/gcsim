import type { Env } from "../bindings";
import { newShareKey } from "./shareKey";

const MAX_ATTEMPTS = 5;

async function gzip(data: ArrayBuffer): Promise<ArrayBuffer> {
	const stream = new Blob([data])
		.stream()
		.pipeThrough(new CompressionStream("gzip"));
	return new Response(stream).arrayBuffer();
}

export async function storeShare(env: Env, data: ArrayBuffer): Promise<string> {
	const gz = await gzip(data);
	for (let i = 0; i < MAX_ATTEMPTS; i++) {
		const key = newShareKey();
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

function gunzip(obj: R2ObjectBody): ReadableStream {
	return obj.body.pipeThrough(new DecompressionStream("gzip"));
}

async function findShare(
	env: Env,
	key: string,
): Promise<R2ObjectBody | Response> {
	const obj = await env.GCSIM_SHARES.get(key);
	return obj ?? fetch(new Request(`${env.API_ENDPOINT}/api/share/${key}`));
}

export async function shareResponse(env: Env, key: string): Promise<Response> {
	const found = await findShare(env, key);
	if (found instanceof Response) {
		const res = new Response(found.body, found);
		res.headers.append("Content-Encoding", "gzip");
		return res;
	}
	// encodeBody "auto" gzips this; a "manual" body is gzipped again by cache.put after clone()
	return new Response(gunzip(found), {
		headers: {
			"Content-Type": "application/json",
			"Content-Encoding": "gzip",
		},
	});
}

export async function decodedShareResponse(
	env: Env,
	key: string,
): Promise<Response> {
	const found = await findShare(env, key);
	if (found instanceof Response) {
		return found;
	}
	return new Response(gunzip(found));
}
