import { env, exports } from "cloudflare:workers";
import { afterEach, describe, expect, it, vi } from "vitest";
import body from "./fixtures/k3.body?raw";
import header from "./fixtures/k3.header?raw";

const TEST_KEY =
	"000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f";
const KEY_RE = /^[6789BCDFGHJKLMNPQRTWbcdfghjkmnpqrtwz]{12}$/;
const MiB = 1024 * 1024;

let nextIP = 0;

// The rate-limiter's window is keyed by IP and persists across tests, so each
// post defaults to a fresh IP; rate-limit tests pass a fixed one to accumulate.
function post(data: BodyInit, auth?: string, ip?: string): Promise<Response> {
	const headers = new Headers();
	if (auth != null) headers.set("X-GCSIM-SHARE-AUTH", auth);
	headers.set(
		"CF-Connecting-IP",
		ip ?? `10.0.${nextIP >> 8}.${nextIP++ & 255}`,
	);
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

	it.each([
		["a missing header", undefined],
		["an empty header", ""],
	])("rejects %s with 403 and writes nothing", async (_, auth) => {
		const res = await post(body, auth);
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
		"rejects a legacy %s header without calling the backend",
		async (id) => {
			const fetchSpy = mockBackend(
				() => new Response("abc123", { status: 202 }),
			);
			const res = await post(body, `${id}:legacyhash`);
			expect(res.status).toBe(403);
			expect(fetchSpy).not.toHaveBeenCalled();
			expect(await storedKeys()).toEqual([]);
		},
	);

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

	it("retries with a new key when the key is taken", async () => {
		const taken = "666666666666";
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

describe("POST /api/share rate limiting", () => {
	const LIMIT = 10;

	it("returns 429 once an IP exceeds the limit", async () => {
		for (let i = 0; i < LIMIT; i++) {
			const res = await post(body, "k3:bad", "1.2.3.4");
			expect(res.status).not.toBe(429);
		}
		const res = await post(body, "k3:bad", "1.2.3.4");
		expect(res.status).toBe(429);
		expect(await storedKeys()).toEqual([]);
	});

	it("counts each IP separately", async () => {
		for (let i = 0; i <= LIMIT; i++) {
			await post(body, "k3:bad", "1.2.3.4");
		}
		const res = await post(body, header, "9.9.9.9");
		expect(res.status).toBe(202);
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

	it("serves the same bytes from the cache on a repeat GET", async () => {
		const key = await (await post(body, header)).text();
		const url = `https://gcsim.test/api/share/${key}`;
		const init = { headers: { "Accept-Encoding": "gzip" } };
		await (await get(`/api/share/${key}`, init)).arrayBuffer();
		await vi.waitFor(async () => {
			expect(await caches.default.match(url)).toBeDefined();
		});
		const res = await get(`/api/share/${key}`, init);
		expect(res.headers.get("Content-Encoding")).toBe("gzip");
		expect(await gunzip(res.body!)).toBe(body);
	});

	it("caches a found share at the edge for a year", async () => {
		const key = await (await post(body, header)).text();
		const res = await get(`/api/share/${key}`);
		expect(res.headers.get("Cache-Control")).toBe(
			"max-age=14400, s-maxage=31536000",
		);
	});

	it.each([404, 500])(
		"does not cache a backend %i for a legacy key",
		async (status) => {
			const key = "V1StGXR8_Z5j";
			const fetchSpy = mockBackend(() => new Response("err", { status }));
			const res = await get(`/api/share/${key}`);
			expect(res.status).toBe(status);
			await res.arrayBuffer();
			await new Promise((r) => setTimeout(r, 50));
			expect(
				await caches.default.match(`https://gcsim.test/api/share/${key}`),
			).toBeUndefined();
			await get(`/api/share/${key}`);
			expect(fetchSpy).toHaveBeenCalledTimes(2);
		},
	);

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
		const res = await get("/api/share/k3j9xqbdHT7w");
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
