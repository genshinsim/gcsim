// One-off generator for the SatoriPreviewCard "Composited" story fixture.
//
// The production OG card composites each portrait (element bg + avatar + weapon
// + artifact set(s), white outline & two-set slice baked in) with Photon inside
// the edge Worker — a Node/wasm step Storybook can't run in the browser. This
// script runs the SAME compositor against real assets (fetched from the live
// asset host) and writes the resulting portraits as PNG files under
// ../src/stories/samples/satoriPortraits/, which the story renders as its
// visual smoke test of the outline + slice.
//
// Run from the storybook package: `npx tsx scripts/genSatoriPortraits.ts`.
// Re-run only when the sample or the compositor changes.
import { mkdirSync, writeFileSync } from "node:fs";
import { fileURLToPath } from "node:url";
import {
	compositePortraits,
	enumerateAssetPaths,
} from "@gcsim/components/src/Cards/SatoriPreviewCard";
import {
	SATORI_SAMPLE_SCALE,
	satoriSample,
} from "../src/stories/samples/satoriSample";

const ASSET_BASE = "https://gcsim.app/api/assets";

async function main() {
	const paths = enumerateAssetPaths(satoriSample);
	const bytes = new Map<string, Uint8Array>();
	await Promise.all(
		paths.map(async (path) => {
			const resp = await fetch(`${ASSET_BASE}/${path}`);
			if (!resp.ok) {
				throw new Error(`fetch ${path}: ${resp.status}`);
			}
			bytes.set(path, new Uint8Array(await resp.arrayBuffer()));
		}),
	);

	const portraits = compositePortraits(
		satoriSample,
		(p) => bytes.get(p),
		SATORI_SAMPLE_SCALE,
	);

	const dir = fileURLToPath(
		new URL("../src/stories/samples/satoriPortraits/", import.meta.url),
	);
	mkdirSync(dir, { recursive: true });
	portraits.forEach((uri, i) => {
		const png = Buffer.from(uri.slice(uri.indexOf(",") + 1), "base64");
		writeFileSync(`${dir}portrait${i}.png`, png);
	});
	console.log(`wrote ${portraits.length} portrait PNGs -> ${dir}`);
}

main().catch((e) => {
	console.error(e);
	process.exit(1);
});
