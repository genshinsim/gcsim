const KEY_CHARS = "6789BCDFGHJKLMNPQRTWbcdfghjkmnpqrtwz";
const KEY_LEN = 12;

export function newShareKey(): string {
	const unbiasedLimit = 256 - (256 % KEY_CHARS.length);
	let out = "";
	while (out.length < KEY_LEN) {
		for (const b of crypto.getRandomValues(new Uint8Array(KEY_LEN))) {
			if (b < unbiasedLimit && out.length < KEY_LEN) {
				out += KEY_CHARS[b % KEY_CHARS.length];
			}
		}
	}
	return out;
}
