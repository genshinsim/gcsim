export type KeyClass = "prod" | "dev";

type ShareKeys = Record<string, { class: KeyClass; key: string }>;

function hexBytes(hex: string): Uint8Array {
	const out = new Uint8Array(hex.length / 2);
	for (let i = 0; i < out.length; i++) {
		out[i] = parseInt(hex.slice(i * 2, i * 2 + 2), 16);
	}
	return out;
}

function base64Bytes(s: string): Uint8Array | null {
	try {
		return Uint8Array.from(atob(s), (c) => c.charCodeAt(0));
	} catch {
		return null;
	}
}

export async function verifyShare(
	keysJSON: string,
	id: string,
	sig: string,
	body: ArrayBuffer,
): Promise<KeyClass | null> {
	const keys: ShareKeys = JSON.parse(keysJSON);
	const entry = new Map(Object.entries(keys)).get(id);
	const tag = base64Bytes(sig);
	if (entry == null || tag == null) {
		return null;
	}
	const key = await crypto.subtle.importKey(
		"raw",
		hexBytes(entry.key),
		{ name: "HMAC", hash: "SHA-256" },
		false,
		["verify"],
	);
	const ok = await crypto.subtle.verify("HMAC", key, tag, body);
	return ok ? entry.class : null;
}
