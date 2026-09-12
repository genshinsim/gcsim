import { readFileSync } from "node:fs";
import { fileURLToPath } from "node:url";

/**
 * The known-good, minimal Sucrose config used by the smoke spec.
 *
 * It is committed as a plain-text fixture (fixtures/sucrose.txt) so it reads
 * exactly like a config a user would paste. `iteration=1` keeps the run fast —
 * the spec asserts structure, never specific numbers.
 */
export const sucroseConfig: string = readFileSync(
	fileURLToPath(new URL("../fixtures/sucrose.txt", import.meta.url)),
	"utf8",
);
