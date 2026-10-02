import { env, exports } from "cloudflare:workers";
import { afterEach, describe, expect, it, vi } from "vitest";
import body from "./fixtures/k3.body?raw";
import header from "./fixtures/k3.header?raw";

const TEST_KEY =
	"000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f";
const SUFFIX = "[6789BCDFGHJKLMNPQRTWbcdfghjkmnpqrtwz]{8}";
const KEY_RE = new RegExp(`^raiden-kazuha-kokomi-furina-${SUFFIX}$`);
const MiB = 1024 * 1024;

function post(data: BodyInit, auth?: string): Promise<Response> {
	const headers = new Headers();
	if (auth != null) headers.set("X-GCSIM-SHARE-AUTH", auth);
	return exports.default.fetch(
		new Request("https://gcsim.test/api/share", {
			method: "POST",
			body: data,
			headers,
		}),
	);
}

function get(path: string, init?: RequestInit): Promise<Response> {
	return exports.default.fetch(new Request(`https://gcsim.test${path}`, init));
}

async function sign(id: string, data: string): Promise<string> {
	const raw = Uint8Array.from(TEST_KEY.match(/../g)!, (h) => parseInt(h, 16));
	const key = await crypto.subtle.importKey(
		"raw",
		raw,
		{ name: "HMAC", hash: "SHA-256" },
		false,
		["sign"],
	);
	const mac = await crypto.subtle.sign(
		"HMAC",
		key,
		new TextEncoder().encode(data),
	);
	return `${id}:${btoa(String.fromCharCode(...new Uint8Array(mac)))}`;
}

function gunzip(data: ReadableStream): Promise<string> {
	return new Response(data.pipeThrough(new DecompressionStream("gzip"))).text();
}

async function storedKeys(): Promise<string[]> {
	return (await env.GCSIM_SHARES.list()).objects.map((o) => o.key);
}

function mockBackend(handler: (req: Request) => Response | Promise<Response>) {
	return vi
		.spyOn(globalThis, "fetch")
		.mockImplementation((input, init) =>
			Promise.resolve(handler(new Request(input, init))),
		);
}

afterEach(async () => {
	vi.restoreAllMocks();
	const keys = await storedKeys();
	if (keys.length > 0) {
		await env.GCSIM_SHARES.delete(keys);
	}
});

describe("POST /api/share", () => {
	it("stores a Go-signed result gzipped under a new key", async () => {
		const res = await post(body, header);
		expect(res.status).toBe(202);
		const key = await res.text();
		expect(key).toMatch(KEY_RE);
		const obj = await env.GCSIM_SHARES.get(key);
		expect(await gunzip(obj!.body)).toBe(body);
	});

	it.each([
		["a body byte", body.replace("raiden char", "raiden chas"), header],
		[
			"a tag byte",
			body,
			header.replace(/.=$/, (m) => (m === "A=" ? "B=" : "A=")),
		],
		["a malformed k3 header", body, "k3:not base64!"],
		["a bare k3 header", body, "k3"],
		["an unknown key id", body, header.replace(/^k3/, "k9")],
	])("rejects %s with 400 and writes nothing", async (_, data, auth) => {
		const res = await post(data, auth);
		expect(res.status).toBe(400);
		expect(await storedKeys()).toEqual([]);
	});

	it("rejects a missing header with 403 and writes nothing", async () => {
		const res = await post(body);
		expect(res.status).toBe(403);
		expect(await storedKeys()).toEqual([]);
	});

	it("rejects a key_type that does not match the key class", async () => {
		const res = await post(body, await sign("k4", body));
		expect(res.status).toBe(400);
		expect(await storedKeys()).toEqual([]);
	});

	it("rejects a signed body that is not JSON", async () => {
		const res = await post("not json", await sign("k3", "not json"));
		expect(res.status).toBe(400);
	});

	it("rejects a signed body without sim_version", async () => {
		const data = JSON.stringify({ key_type: "prod" });
		const res = await post(data, await sign("k3", data));
		expect(res.status).toBe(400);
	});

	it.each(["prod", "dev"])(
		"forwards a %s header to the backend unchanged",
		async (id) => {
			const auth = `${id}:legacyhash`;
			const fetchSpy = mockBackend(
				() => new Response("abc123", { status: 202 }),
			);
			const res = await post(body, auth);
			expect(res.status).toBe(202);
			expect(await res.text()).toBe("abc123");
			const sent = new Request(
				...(fetchSpy.mock.calls[0] as [RequestInfo, RequestInit]),
			);
			expect(sent.url).toBe("https://backend.test/api/share");
			expect(sent.method).toBe("POST");
			expect(sent.headers.get("X-GCSIM-SHARE-AUTH")).toBe(auth);
			expect(await sent.text()).toBe(body);
			expect(await storedKeys()).toEqual([]);
		},
	);

	it("passes the backend's error through for a legacy header", async () => {
		mockBackend(() => new Response("Bad Request", { status: 400 }));
		const res = await post(body, "prod:legacyhash");
		expect(res.status).toBe(400);
	});

	it("accepts exactly 1 MiB", async () => {
		const data = padTo(MiB);
		const res = await post(data, await sign("k3", data));
		expect(res.status).toBe(202);
	});

	it("rejects 1 MiB + 1 byte with 413 and writes nothing", async () => {
		const data = padTo(MiB + 1);
		const res = await post(data, await sign("k3", data));
		expect(res.status).toBe(413);
		expect(await storedKeys()).toEqual([]);
	});

	it("retries with a new suffix when the key is taken", async () => {
		const taken = "raiden-kazuha-kokomi-furina-66666666";
		await env.GCSIM_SHARES.put(taken, "existing");
		const real = crypto.getRandomValues.bind(crypto);
		vi.spyOn(crypto, "getRandomValues")
			.mockImplementationOnce((a) => {
				new Uint8Array(a.buffer).fill(0);
				return a;
			})
			.mockImplementation(real);
		const res = await post(body, header);
		expect(res.status).toBe(202);
		const key = await res.text();
		expect(key).not.toBe(taken);
		expect(key).toMatch(KEY_RE);
		expect(await (await env.GCSIM_SHARES.get(taken))!.text()).toBe("existing");
	});
});

describe("GET /api/share/:key", () => {
	it("serves a stored share as gzip that decodes once to the posted bytes", async () => {
		const key = await (await post(body, header)).text();
		const fetchSpy = mockBackend(() => new Response(null, { status: 500 }));
		const res = await get(`/api/share/${key}`, {
			headers: { "Accept-Encoding": "gzip" },
		});
		expect(res.status).toBe(200);
		expect(res.headers.get("Content-Type")).toBe("application/json");
		expect(res.headers.get("Content-Encoding")).toBe("gzip");
		const text = await res.text();
		expect(text).toBe(body);
		expect(JSON.parse(text).sim_version).toBe("fixture");
		expect(fetchSpy).not.toHaveBeenCalled();
	});

	it.each(["V1StGXR8_Z5j", "0b5bd9a6-8a3e-4b3c-9f1e-2d4c6b8a0e1f"])(
		"falls back to the backend for legacy key %s",
		async (key) => {
			const fetchSpy = mockBackend(
				() => new Response('{"sim_version":"legacy"}', { status: 200 }),
			);
			const res = await get(`/api/share/${key}`);
			expect(res.status).toBe(200);
			expect(new Request(fetchSpy.mock.calls[0][0]).url).toBe(
				`https://backend.test/api/share/${key}`,
			);
		},
	);

	it("returns 404 for an unknown new-format key", async () => {
		mockBackend(() => new Response("not found", { status: 404 }));
		const res = await get("/api/share/raiden-k3j9xqbd");
		expect(res.status).toBe(404);
	});

	it("no longer serves /api/share/random", async () => {
		const fetchSpy = mockBackend(() => new Response("{}", { status: 200 }));
		const res = await get("/api/share/random");
		expect(res.status).toBe(404);
		expect(fetchSpy).not.toHaveBeenCalled();
	});
});

function padTo(size: number): string {
	const base = JSON.stringify({ sim_version: "x", key_type: "prod", pad: "" });
	return base.replace('"pad":""', `"pad":"${"a".repeat(size - base.length)}"`);
}

describe("share pages", () => {
	it("renders the OG preview card for a stored share", async () => {
		const key = await (await post(body, header)).text();
		const fetchSpy = mockBackend(() => new Response(null, { status: 404 }));
		const res = await get(`/api/preview/${key}.png`);
		expect(res.status).toBe(200);
		expect(res.headers.get("Content-Type")).toBe("image/png");
		const urls = fetchSpy.mock.calls.map((c) => new Request(c[0]).url);
		expect(urls).not.toContain(`https://backend.test/api/share/${key}`);
	});

	it("injects the OG head for a stored share", async () => {
		const key = await (await post(body, header)).text();
		const res = await get(`/sh/${key}`);
		expect(res.status).toBe(200);
		expect(await res.text()).toContain(
			`<meta property="og:image" content="https://gcsim.test/api/preview/${key}.png" />`,
		);
	});
});
