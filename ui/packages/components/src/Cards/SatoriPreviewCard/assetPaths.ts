import type { model } from "@gcsim/types";

// Asset paths the card requests, relative to the asset base (e.g.
// "avatar/Nahida.png"). Centralised so the card's <img> srcs and any external
// pre-fetcher (the edge worker) agree on the exact same set of paths.

/** Injectable seam: maps a relative asset path to an <img src>. */
export type ResolveAsset = (path: string) => string;

export const NAHIDA_PLACEHOLDER_PATH = "misc/nahida.png";
// Names come from the (all-optional) model; interpolate them as the card
// always has, so paths stay byte-identical to the pre-seam <img src>s.
export const avatarPath = (name: string | undefined): string =>
	`avatar/${name}.png`;
export const weaponPath = (name: string | undefined): string =>
	`weapons/${name}.png`;
export const artifactFlowerPath = (set: string): string =>
	`artifacts/${set}_flower.png`;

// Enumerate the deduped asset paths <SatoriPreviewCard> requests for a result,
// mirroring Portraits' render logic exactly. Lets a caller pre-fetch those
// bytes and feed them back through resolveAsset so Satori makes no image
// fetches of its own.
export function enumerateAssetPaths(data: model.SimulationResult): string[] {
	const paths = new Set<string>();
	for (const char of data.character_details ?? []) {
		// Empty slot renders the Nahida placeholder. Guard matches Portrait's
		// strict `char === null` exactly, so the two never disagree on which
		// paths a slot needs (a mismatch would cache a fallback render).
		if (char === null) {
			paths.add(NAHIDA_PLACEHOLDER_PATH);
			continue;
		}
		paths.add(avatarPath(char.name));
		if (char.weapon?.name) {
			paths.add(weaponPath(char.weapon.name));
		}
		const sets = char.sets ? Object.keys(char.sets) : [];
		if (sets.length > 0) {
			paths.add(artifactFlowerPath(sets[0]));
			if (sets.length > 1) {
				paths.add(artifactFlowerPath(sets[1]));
			}
		}
	}
	return [...paths];
}
