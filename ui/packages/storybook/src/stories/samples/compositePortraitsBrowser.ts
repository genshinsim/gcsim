// Runs the REAL portrait compositor in the browser, so the CompositedLive story
// reflects edits to portraitCompositor.ts on hot-reload — the way to eyeball
// outline/slice changes without re-running the Node fixture generator.
//
// DEV ONLY: asset bytes are fetched through Storybook's dev `/api` proxy (→
// gcsim.app, see .storybook/middleware.js), which a static `build-storybook` /
// Chromatic build doesn't have. Those builds use the baked `Composited` story.
import { initPhoton } from "@cf-wasm/photon/others";
// Vite serves the wasm as a URL; hand it to Photon's async init (the `others`
// build, unlike `workerd`/`node`, does not self-initialise).
import photonWasmUrl from "@cf-wasm/photon/photon.wasm?url";
import { enumerateAssetPaths } from "@gcsim/components/src/Cards/SatoriPreviewCard/assetPaths";
import { compositePortraits } from "@gcsim/components/src/Cards/SatoriPreviewCard/portraitCompositor";
import type { model } from "@gcsim/types";
import { SATORI_SAMPLE_SCALE } from "./satoriSample";

// initPhoton throws if called twice; across HMR the photon module stays cached,
// so init once and await the existing promise thereafter.
async function ensurePhoton(): Promise<void> {
	if (initPhoton.initialized) {
		await initPhoton.ensure();
	} else {
		await initPhoton(photonWasmUrl);
	}
}

export async function compositePortraitsInBrowser(
	data: model.SimulationResult,
): Promise<string[]> {
	await ensurePhoton();
	const bytes = new Map<string, Uint8Array>();
	await Promise.all(
		enumerateAssetPaths(data).map(async (path) => {
			const resp = await fetch(`/api/assets/${path}`);
			if (resp.ok) {
				bytes.set(path, new Uint8Array(await resp.arrayBuffer()));
			}
		}),
	);
	return compositePortraits(data, (p) => bytes.get(p), SATORI_SAMPLE_SCALE);
}
