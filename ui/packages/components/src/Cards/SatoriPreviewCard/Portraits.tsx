import type { model } from "@gcsim/types";
import {
	artifactFlowerPath,
	avatarPath,
	NAHIDA_PLACEHOLDER_PATH,
	type ResolveAsset,
	weaponPath,
} from "./assetPaths";
import { elementBackgrounds } from "./backgrounds";
import { GRAY_400, GRAY_600, GRAY_700, levelColor, PRIMARY_BG } from "./colors";
import { FONT_FAMILY } from "./fonts";
// Portrait layout constants (logical px), shared with the Photon compositor so
// the layered and composited renders place every layer identically.
import {
	ARTIFACT_BOTTOM,
	ARTIFACT_HALF,
	ARTIFACT_LEFT,
	ARTIFACT_SIZE,
	AVATAR_MARGIN_TOP,
	AVATAR_SIZE,
	ICON_OPACITY,
	PLACEHOLDER_OPACITY,
	WEAPON_BOTTOM,
	WEAPON_RIGHT,
	WEAPON_SIZE,
	WIP_TOP,
} from "./portraitGeometry";

// Live badge accent colors (text-geo / text-electro), read from the rendered
// PreviewCard.
const CONS_COLOR = "#f8ba4e"; // text-geo
const REFINE_COLOR = "#b25dcd"; // text-electro

function portraitBackground(element: string): string {
	return elementBackgrounds[element] ?? elementBackgrounds.default;
}

type PortraitProps = {
	char: model.Character | null;
	i: number;
	invalid: boolean;
	width: number;
	height: number;
	margin: number;
	resolveAsset: ResolveAsset;
	// Pre-composited portrait image (bg + avatar + weapon + artifacts, with the
	// white outline and two-set slice baked in). When set, the imagery renders as
	// a single <img> and only the text badges are layered on top; when absent, the
	// portrait is stacked layer-by-layer through resolveAsset (browser/Storybook).
	composited?: string;
};

// A single character portrait: avatar, weapon, artifact set(s) and level/cons
// badges. Imagery <img> srcs come from the injected resolveAsset seam; the
// element background is a pre-blended bundled image (see ./backgrounds).
const Portrait = ({
	char,
	i,
	invalid,
	width,
	height,
	margin,
	resolveAsset,
	composited,
}: PortraitProps) => {
	const base = {
		display: "flex" as const,
		width,
		height,
		margin,
		borderRadius: 4,
		border: `1px solid ${GRAY_600}`,
		overflow: "hidden" as const,
		position: "relative" as const,
	};

	// Empty slot.
	if (char === null) {
		// Composited: the placeholder is already baked into the portrait image.
		if (composited) {
			return (
				<div style={base}>
					<img src={composited} width={width} height={height} alt="" />
				</div>
			);
		}
		return (
			<div
				style={{
					...base,
					backgroundColor: GRAY_400,
					alignItems: "center",
					justifyContent: "center",
				}}
			>
				<img
					src={resolveAsset(NAHIDA_PLACEHOLDER_PATH)}
					width={AVATAR_SIZE}
					height={AVATAR_SIZE}
					alt=""
					style={{ objectFit: "contain", opacity: PLACEHOLDER_OPACITY }}
				/>
			</div>
		);
	}

	const sets = char.sets ? Object.keys(char.sets) : [];
	// A lone 2-piece set renders as a half-width flower (matches AvatarPortrait).
	const isHalfWidthSet = sets.length === 1 && char.sets?.[sets[0]] === 2;
	const twoSets = sets.length > 1;

	// cons / refine and level badges (font-mono, text-xs, gray-700 @ 85%).
	const badgeStyle = {
		display: "flex" as const,
		flexDirection: "row" as const,
		gap: 4,
		padding: "0 4px",
		backgroundColor: GRAY_700,
		opacity: 0.85,
		fontFamily: FONT_FAMILY,
		fontSize: 12,
		fontWeight: 700,
	};

	// Text badges, shared by the layered and composited renders.
	const consBadge = (
		<div
			style={{
				...badgeStyle,
				position: "absolute",
				left: 0,
				top: 0,
				borderTopLeftRadius: 4,
				borderBottomRightRadius: 8,
			}}
		>
			<span style={{ color: CONS_COLOR }}>{`C${char.cons ?? 0}`}</span>
			{char.weapon ? (
				<span
					style={{ color: REFINE_COLOR }}
				>{`R${char.weapon.refine ?? 0}`}</span>
			) : null}
		</div>
	);

	const levelBadge = (
		<div
			style={{
				...badgeStyle,
				position: "absolute",
				right: 0,
				top: 0,
				alignItems: "center",
				borderTopRightRadius: 4,
				borderBottomLeftRadius: 8,
			}}
		>
			<span style={{ color: GRAY_400, fontWeight: 400 }}>lvl</span>
			<span style={{ color: levelColor(i) }}>{char.level}</span>
		</div>
	);

	// incomplete build: full-width WIP bar across the portrait (top-1/3).
	const wipOverlay = invalid ? (
		<div
			style={{
				position: "absolute",
				top: WIP_TOP,
				left: 0,
				width,
				display: "flex",
				alignItems: "center",
				justifyContent: "center",
				padding: "6px 0",
				border: "1px solid transparent",
				backgroundColor: PRIMARY_BG,
			}}
		>
			<span
				style={{
					color: "#ef4444",
					fontFamily: FONT_FAMILY,
					fontWeight: 700,
					fontSize: 12,
					lineHeight: "16px",
					textTransform: "uppercase",
				}}
			>
				WIP
			</span>
		</div>
	) : null;

	// Composited: one flat image for all imagery, text badges layered on top.
	if (composited) {
		return (
			<div style={base}>
				<img
					src={composited}
					width={width}
					height={height}
					alt={char.name ?? ""}
				/>
				{consBadge}
				{levelBadge}
				{wipOverlay}
			</div>
		);
	}

	return (
		<div
			style={{
				...base,
				alignItems: "flex-start",
				justifyContent: "center",
				backgroundImage: `url(${portraitBackground(char.element ?? "")})`,
				// The bundled backgrounds share the portrait's aspect ratio
				// (381x318 ~= 127x106), so an exact 100% 100% fill matches the
				// live card's `cover` with no visible distortion. `cover` +
				// `center` is avoided on purpose: Satori's pattern-position math
				// emits x="NaN" y="NaN" for it, which makes the raster tile the
				// background from 0,0 instead of filling the portrait.
				backgroundSize: "100% 100%",
				backgroundRepeat: "no-repeat",
			}}
		>
			{/* avatar (h-24, top-aligned under the card's pt-2) */}
			<img
				src={resolveAsset(avatarPath(char.name))}
				width={AVATAR_SIZE}
				height={AVATAR_SIZE}
				alt={char.name ?? ""}
				style={{ objectFit: "contain", marginTop: AVATAR_MARGIN_TOP }}
			/>

			{/* weapon */}
			{char.weapon?.name ? (
				<img
					src={resolveAsset(weaponPath(char.weapon.name))}
					width={WEAPON_SIZE}
					height={WEAPON_SIZE}
					alt=""
					style={{
						position: "absolute",
						bottom: WEAPON_BOTTOM,
						right: WEAPON_RIGHT,
						objectFit: "contain",
						opacity: ICON_OPACITY,
					}}
				/>
			) : null}

			{/* artifact set(s) — full flower, or two half-width halves for a 2-set build */}
			{sets.length > 0 ? (
				<div
					style={{
						position: "absolute",
						bottom: ARTIFACT_BOTTOM,
						left: ARTIFACT_LEFT,
						display: "flex",
						flexDirection: "row",
						opacity: ICON_OPACITY,
					}}
				>
					<img
						src={resolveAsset(artifactFlowerPath(sets[0]))}
						width={twoSets || isHalfWidthSet ? ARTIFACT_HALF : ARTIFACT_SIZE}
						height={ARTIFACT_SIZE}
						alt=""
						style={{ objectFit: "cover" }}
					/>
					{twoSets ? (
						<img
							src={resolveAsset(artifactFlowerPath(sets[1]))}
							width={ARTIFACT_HALF}
							height={ARTIFACT_SIZE}
							alt=""
							style={{ objectFit: "cover" }}
						/>
					) : null}
				</div>
			) : null}

			{consBadge}
			{levelBadge}
			{wipOverlay}
		</div>
	);
};

type Props = {
	data: model.SimulationResult;
	width: number;
	height: number;
	margin: number;
	resolveAsset: ResolveAsset;
	// Pre-composited portrait images, one per character slot (see Portrait).
	composited?: readonly (string | null)[];
};

// The four-portrait row (flexbox, replacing the live card's CSS grid). Each
// portrait carries its own 4px margin, reproducing grid-cols-4 + m-1.
export const Portraits = ({
	data,
	width,
	height,
	margin,
	resolveAsset,
	composited,
}: Props) => {
	const chars = data.character_details ?? [];
	return (
		<div
			style={{
				display: "flex",
				flexDirection: "row",
				justifyContent: "center",
			}}
		>
			{chars.map((c, i) => (
				<Portrait
					key={c.name ?? `empty-${i}`}
					char={c}
					i={i}
					invalid={data.incomplete_characters?.includes(c.name ?? "") ?? false}
					width={width}
					height={height}
					margin={margin}
					resolveAsset={resolveAsset}
					composited={composited?.[i] ?? undefined}
				/>
			))}
		</div>
	);
};
