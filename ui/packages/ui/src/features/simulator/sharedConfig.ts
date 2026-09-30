import Pako from "pako";

export type SharedConfig = {
	config: string;
	title?: string;
	source?: string;
};

function isNonEmptyString(v: unknown): v is string {
	return typeof v === "string" && v !== "";
}

// Accepts url-safe and standard base64; spaces are '+' mangled by query-string decoding.
function base64ToBytes(encoded: string): Uint8Array {
	const b64 = encoded
		.replaceAll(" ", "+")
		.replaceAll("-", "+")
		.replaceAll("_", "/");
	const binary = atob(b64.padEnd(Math.ceil(b64.length / 4) * 4, "="));
	return Uint8Array.from(binary, (c) => c.charCodeAt(0));
}

function bytesToBase64(bytes: Uint8Array): string {
	let binary = "";
	for (const b of bytes) binary += String.fromCharCode(b);
	return btoa(binary);
}

export function decodeSharedConfig(encoded: string): SharedConfig {
	const json = Pako.inflate(base64ToBytes(encoded), { to: "string" });
	const raw: unknown = JSON.parse(json);
	if (raw === null || typeof raw !== "object") {
		throw new Error("shared config is not a JSON object");
	}
	const { config, title, source } = raw as Record<string, unknown>;
	if (!isNonEmptyString(config)) {
		throw new Error("shared config is missing a config");
	}
	return {
		config,
		...(isNonEmptyString(title) && { title }),
		...(isNonEmptyString(source) && { source }),
	};
}

export function encodeSharedConfig(payload: SharedConfig): string {
	return bytesToBase64(Pako.gzip(JSON.stringify(payload)));
}
