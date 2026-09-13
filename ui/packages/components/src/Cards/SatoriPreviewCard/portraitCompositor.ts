// Photon (WASM) portrait compositor for the Satori OG preview.
//
// Satori cannot emit the SVG filters the live AvatarPortrait uses (a hard white
// `feMorphology` dilate outline) nor honour `preserveAspectRatio` slicing, so
// the small weapon/artifact icons rendered straight through Satori wash out and
// two-set builds show two identical centre crops. To restore the old look we
// pre-composite each character's portrait — element background, avatar, weapon
// and artifact set(s) — into ONE flat PNG here, with the white outline and the
// two-set half/half slice baked in at raster level. The Satori card then draws
// a single <img> per slot and keeps only the text badges as Satori elements.
//
// Runs on the edge Worker (import resolves to `@cf-wasm/photon`'s `workerd`
// build) and in a Node harness (`node` build) — the API is identical, so the
// same module drives both. It is intentionally NOT re-exported from the browser
// `Cards` barrel: nothing in the web/Storybook bundle should pull in the wasm.
import {
	crop,
	PhotonImage,
	resize,
	SamplingFilter,
	watermark,
} from "@cf-wasm/photon";
import type { model } from "@gcsim/types";
import {
	artifactFlowerPath,
	avatarPath,
	NAHIDA_PLACEHOLDER_PATH,
	weaponPath,
} from "./assetPaths";
import { elementBackgrounds } from "./backgrounds";
import { GRAY_400 } from "./colors";
// Portrait layout constants (logical px), shared with the layered Portrait
// render so the two agree. Everything below multiplies by the render scale
// before touching pixels.
import {
	ARTIFACT_BOTTOM,
	ARTIFACT_LEFT,
	ARTIFACT_SIZE,
	AVATAR_MARGIN_TOP,
	AVATAR_SIZE,
	ICON_OPACITY,
	PLACEHOLDER_OPACITY,
	PORTRAIT_H,
	PORTRAIT_W,
	WEAPON_BOTTOM,
	WEAPON_RIGHT,
	WEAPON_SIZE,
} from "./portraitGeometry";

// Resolve a relative asset path (see assetPaths) to raw image bytes. Mirrors the
// Worker's in-process resolver output; a miss yields the placeholder bytes.
export type ResolveBytes = (path: string) => Uint8Array | undefined;

// Empty-slot background (tailwind gray-400), matching Portrait's char===null case.
const EMPTY_BG = hexToRgb(GRAY_400);

function hexToRgb(hex: string): [number, number, number] {
	const n = Number.parseInt(hex.slice(1), 16);
	return [(n >> 16) & 0xff, (n >> 8) & 0xff, n & 0xff];
}

function even(n: number): number {
	// Force even so the two-set slice splits on an integer boundary.
	return Math.round(n / 2) * 2;
}

function dataUriToBytes(uri: string): Uint8Array {
	const base64 = uri.slice(uri.indexOf(",") + 1);
	const binary = atob(base64);
	const out = new Uint8Array(binary.length);
	for (let i = 0; i < binary.length; i++) {
		out[i] = binary.charCodeAt(i);
	}
	return out;
}

function decode(bytes: Uint8Array): PhotonImage {
	return PhotonImage.new_from_byteslice(bytes);
}

function resizeTo(img: PhotonImage, w: number, h: number): PhotonImage {
	const out = resize(img, w, h, SamplingFilter.Lanczos3);
	img.free();
	return out;
}

// Scale every pixel's alpha by `factor` (0..1), for the live opacity-85/50 wraps.
function withOpacity(img: PhotonImage, factor: number): PhotonImage {
	if (factor >= 1) {
		return img;
	}
	const px = img.get_raw_pixels();
	for (let i = 3; i < px.length; i += 4) {
		px[i] = Math.round(px[i] * factor);
	}
	const out = new PhotonImage(px, img.get_width(), img.get_height());
	img.free();
	return out;
}

// A hard white edge hugging the icon's alpha silhouette (Photon has no dilate).
// Pads by R so the halo isn't clipped, stamps a white silhouette at every offset
// within Chebyshev radius R (a box structuring element, matching the old
// feMorphology dilate), then draws the icon on top (source-over).
function withWhiteOutline(img: PhotonImage, r: number): PhotonImage {
	const w = img.get_width();
	const h = img.get_height();
	const src = img.get_raw_pixels();
	const nw = w + 2 * r;
	const nh = h + 2 * r;
	const out = new Uint8Array(nw * nh * 4); // zero = transparent

	// Pass 1: white halo. For each opaque source pixel, paint white at every
	// offset within radius R, keeping the strongest alpha seen.
	for (let y = 0; y < h; y++) {
		for (let x = 0; x < w; x++) {
			const a = src[(y * w + x) * 4 + 3];
			// Skip near-transparent antialias fringe (alpha < ~6%) so the outline
			// traces the solid silhouette, not the icon's soft anti-aliased edge.
			if (a < 16) continue;
			for (let dy = -r; dy <= r; dy++) {
				for (let dx = -r; dx <= r; dx++) {
					const di = ((y + r + dy) * nw + (x + r + dx)) * 4;
					if (a > out[di + 3]) {
						out[di] = 255;
						out[di + 1] = 255;
						out[di + 2] = 255;
						out[di + 3] = a;
					}
				}
			}
		}
	}

	// Pass 2: draw the original icon over the halo (source-over), centred in pad.
	for (let y = 0; y < h; y++) {
		for (let x = 0; x < w; x++) {
			const si = (y * w + x) * 4;
			const a = src[si + 3] / 255;
			if (a === 0) continue;
			const di = ((y + r) * nw + (x + r)) * 4;
			const ba = out[di + 3] / 255;
			const oa = a + ba * (1 - a);
			for (let c = 0; c < 3; c++) {
				out[di + c] = Math.round(
					(src[si + c] * a + out[di + c] * ba * (1 - a)) / (oa || 1),
				);
			}
			out[di + 3] = Math.round(oa * 255);
		}
	}
	img.free();
	return new PhotonImage(out, nw, nh);
}

// Build the artifact gear icon at composited resolution, then outline once over
// the combined silhouette. Matches AvatarPortrait/ArtifactsIcon:
//   - two sets: left half of A joined to the right half of B (one clean outline)
//   - lone 2pc set: left-half slice
//   - single non-2pc set: full square flower
function buildGearIcon(
	setImgs: PhotonImage[],
	isLoneHalf: boolean,
	cell: number,
	r: number,
): PhotonImage {
	const half = cell / 2;
	if (setImgs.length >= 2) {
		const left = crop(setImgs[0], 0, 0, half, cell); // left half of A
		const right = crop(setImgs[1], half, 0, cell, cell); // right half of B
		const base = new PhotonImage(new Uint8Array(cell * cell * 4), cell, cell);
		watermark(base, left, 0n, 0n);
		watermark(base, right, BigInt(half), 0n);
		for (const s of setImgs) {
			s.free();
		}
		left.free();
		right.free();
		return withWhiteOutline(base, r);
	}
	if (isLoneHalf) {
		const left = crop(setImgs[0], 0, 0, half, cell);
		const base = new PhotonImage(new Uint8Array(half * cell * 4), half, cell);
		watermark(base, left, 0n, 0n);
		setImgs[0].free();
		left.free();
		return withWhiteOutline(base, r);
	}
	return withWhiteOutline(setImgs[0], r); // single non-2pc → full flower
}

// Draw `overlay` onto `base` with its top-left at (x, y); watermark takes BigInt.
function drawAt(base: PhotonImage, overlay: PhotonImage, x: number, y: number) {
	watermark(base, overlay, BigInt(Math.round(x)), BigInt(Math.round(y)));
}

// Solid opaque RGBA canvas.
function solid(
	w: number,
	h: number,
	[r, g, b]: [number, number, number],
): PhotonImage {
	const px = new Uint8Array(w * h * 4);
	for (let i = 0; i < px.length; i += 4) {
		px[i] = r;
		px[i + 1] = g;
		px[i + 2] = b;
		px[i + 3] = 255;
	}
	return new PhotonImage(px, w, h);
}

function backgroundImage(element: string): PhotonImage {
	const uri = elementBackgrounds[element] ?? elementBackgrounds.default;
	return decode(dataUriToBytes(uri));
}

// Composite one character portrait into a flat PhotonImage at PORTRAIT * scale.
// Never throws: a decode/resize/watermark failure (corrupt bytes, bad dims)
// degrades to the element background — or the empty gray field — so a single
// bad asset can't fail the whole OG render (the old layered path just fetched a
// placeholder). Guarantees compositePortraits always yields a valid image.
function compositePortrait(
	char: model.Character | null,
	resolveBytes: ResolveBytes,
	scale: number,
): PhotonImage {
	const pw = Math.round(PORTRAIT_W * scale);
	const ph = Math.round(PORTRAIT_H * scale);
	try {
		return composePortraitLayers(char, resolveBytes, scale, pw, ph);
	} catch (e) {
		console.error("portrait compositing failed; using plain background", e);
		try {
			return char
				? resizeTo(backgroundImage(char.element ?? ""), pw, ph)
				: solid(pw, ph, EMPTY_BG);
		} catch {
			return solid(pw, ph, EMPTY_BG);
		}
	}
}

// Stack a portrait's layers (element bg + avatar + weapon + artifact set(s), with
// the white outline and two-set slice baked in). May throw on a Photon failure;
// compositePortrait wraps it with the plain-background fallback.
function composePortraitLayers(
	char: model.Character | null,
	resolveBytes: ResolveBytes,
	scale: number,
	pw: number,
	ph: number,
): PhotonImage {
	const r = Math.max(1, Math.round(scale)); // outline radius ~ old 1px dilate

	// Empty slot: gray field + centred, half-opacity Nahida placeholder.
	if (char === null) {
		const canvas = solid(pw, ph, EMPTY_BG);
		const bytes = resolveBytes(NAHIDA_PLACEHOLDER_PATH);
		if (bytes) {
			const size = Math.round(AVATAR_SIZE * scale);
			const nahida = withOpacity(
				resizeTo(decode(bytes), size, size),
				PLACEHOLDER_OPACITY,
			);
			drawAt(canvas, nahida, (pw - size) / 2, (ph - size) / 2);
			nahida.free();
		}
		return canvas;
	}

	// Background fills the whole portrait (Portrait uses backgroundSize 100% 100%).
	const canvas = resizeTo(backgroundImage(char.element ?? ""), pw, ph);

	// Avatar: 96x96, horizontally centred, marginTop 8 (contain; source is square).
	const avatarBytes = resolveBytes(avatarPath(char.name));
	if (avatarBytes) {
		const size = Math.round(AVATAR_SIZE * scale);
		const avatar = resizeTo(decode(avatarBytes), size, size);
		drawAt(canvas, avatar, (pw - size) / 2, AVATAR_MARGIN_TOP * scale);
		avatar.free();
	}

	// Weapon: 55x55, bottom-right (right:-4 overhangs and is clipped), outlined.
	const weaponBytes = char.weapon?.name
		? resolveBytes(weaponPath(char.weapon.name))
		: undefined;
	if (weaponBytes) {
		const size = Math.round(WEAPON_SIZE * scale);
		let weapon = resizeTo(decode(weaponBytes), size, size);
		weapon = withOpacity(withWhiteOutline(weapon, r), ICON_OPACITY);
		// left = W - right - size (CSS `right`): with right:-4 the icon's left is
		// 127 - (-4) - 55 = 76, overhanging the right edge by 4px (canvas clips it).
		const x = (PORTRAIT_W - WEAPON_RIGHT - WEAPON_SIZE) * scale;
		const y = (PORTRAIT_H - WEAPON_SIZE - WEAPON_BOTTOM) * scale;
		drawAt(canvas, weapon, x - r, y - r); // -r: undo the outline pad
		weapon.free();
	}

	// Artifact set(s): 35x35 (or the two-set/lone-2pc slice), bottom-left, outlined.
	const sets = char.sets ? Object.keys(char.sets) : [];
	if (sets.length > 0) {
		const cell = even(Math.round(ARTIFACT_SIZE * scale));
		const isLoneHalf = sets.length === 1 && char.sets?.[sets[0]] === 2;
		const setImgs: PhotonImage[] = [];
		for (const set of sets.slice(0, 2)) {
			const bytes = resolveBytes(artifactFlowerPath(set));
			if (!bytes) {
				break;
			}
			setImgs.push(resizeTo(decode(bytes), cell, cell));
		}
		if (setImgs.length > 0) {
			const gear = withOpacity(
				buildGearIcon(setImgs, isLoneHalf, cell, r),
				ICON_OPACITY,
			);
			const x = ARTIFACT_LEFT * scale;
			const y = (PORTRAIT_H - ARTIFACT_SIZE - ARTIFACT_BOTTOM) * scale;
			drawAt(canvas, gear, x - r, y - r);
			gear.free();
		}
	}

	return canvas;
}

// Composite every character slot's portrait into a flat PNG `data:` URI, in
// character_details order. Feed the result to <SatoriPreviewCard portraits={...}>.
export function compositePortraits(
	data: model.SimulationResult,
	resolveBytes: ResolveBytes,
	scale: number,
): string[] {
	const chars = data.character_details ?? [];
	return chars.map((char) => {
		const img = compositePortrait(char, resolveBytes, scale);
		const b64 = img.get_base64();
		img.free();
		return b64.startsWith("data:") ? b64 : `data:image/png;base64,${b64}`;
	});
}
