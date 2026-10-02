const SUFFIX_CHARS = "6789BCDFGHJKLMNPQRTWbcdfghjkmnpqrtwz";
const SUFFIX_LEN = 8;
const MAX_NAMES = 4;

export function shareKeyPrefix(names: string[]): string {
	return names
		.slice(0, MAX_NAMES)
		.map((n) => n.toLowerCase().replace(/[^a-z0-9]/g, ""))
		.filter((n) => n !== "")
		.join("-");
}

function shareKeySuffix(): string {
	const unbiasedLimit = 256 - (256 % SUFFIX_CHARS.length);
	let out = "";
	while (out.length < SUFFIX_LEN) {
		for (const b of crypto.getRandomValues(new Uint8Array(SUFFIX_LEN))) {
			if (b < unbiasedLimit && out.length < SUFFIX_LEN) {
				out += SUFFIX_CHARS[b % SUFFIX_CHARS.length];
			}
		}
	}
	return out;
}

export function newShareKey(prefix: string): string {
	const suffix = shareKeySuffix();
	return prefix === "" ? suffix : `${prefix}-${suffix}`;
}
