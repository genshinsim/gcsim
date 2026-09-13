import type { model } from "@gcsim/types";
import { elementBackgrounds } from "./backgrounds";
import { GRAY_400, GRAY_600, GRAY_700, levelColor, PRIMARY_BG } from "./colors";
import { FONT_FAMILY } from "./fonts";

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
	assetBase: string;
};

// A single character portrait: avatar, weapon, artifact set(s) and level/cons
// badges. All imagery is remote URL <img> behind the injected asset base; the
// element background is a pre-blended bundled image (see ./backgrounds).
const Portrait = ({
	char,
	i,
	invalid,
	width,
	height,
	margin,
	assetBase,
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
		return (
			<div
				style={{
					...base,
					backgroundColor: "#9ca3af",
					alignItems: "center",
					justifyContent: "center",
				}}
			>
				<img
					src={`${assetBase}/misc/nahida.png`}
					width={96}
					height={96}
					alt=""
					style={{ objectFit: "contain", opacity: 0.5 }}
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

	return (
		<div
			style={{
				...base,
				alignItems: "flex-start",
				justifyContent: "center",
				backgroundImage: `url(${portraitBackground(char.element ?? "")})`,
				backgroundSize: "cover",
				backgroundPosition: "center",
			}}
		>
			{/* avatar (h-24, top-aligned under the card's pt-2) */}
			<img
				src={`${assetBase}/avatar/${char.name}.png`}
				width={96}
				height={96}
				alt={char.name ?? ""}
				style={{ objectFit: "contain", marginTop: 8 }}
			/>

			{/* weapon */}
			{char.weapon?.name ? (
				<img
					src={`${assetBase}/weapons/${char.weapon.name}.png`}
					width={55}
					height={55}
					alt=""
					style={{
						position: "absolute",
						bottom: 4,
						right: -4,
						objectFit: "contain",
						opacity: 0.85,
					}}
				/>
			) : null}

			{/* artifact set(s) — 35x35, or two 17.5-wide halves for a 2-set build */}
			{sets.length > 0 ? (
				<div
					style={{
						position: "absolute",
						bottom: 1,
						left: 1,
						display: "flex",
						flexDirection: "row",
						opacity: 0.85,
					}}
				>
					<img
						src={`${assetBase}/artifacts/${sets[0]}_flower.png`}
						width={twoSets || isHalfWidthSet ? 17.5 : 35}
						height={35}
						alt=""
						style={{ objectFit: "cover" }}
					/>
					{twoSets ? (
						<img
							src={`${assetBase}/artifacts/${sets[1]}_flower.png`}
							width={17.5}
							height={35}
							alt=""
							style={{ objectFit: "cover" }}
						/>
					) : null}
				</div>
			) : null}

			{/* cons / refine (top-left) */}
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

			{/* level (top-right) */}
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

			{/* incomplete build: full-width WIP bar across the portrait (top-1/3) */}
			{invalid ? (
				<div
					style={{
						position: "absolute",
						top: "33%",
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
			) : null}
		</div>
	);
};

type Props = {
	data: model.SimulationResult;
	width: number;
	height: number;
	margin: number;
	assetBase: string;
};

// The four-portrait row (flexbox, replacing the live card's CSS grid). Each
// portrait carries its own 4px margin, reproducing grid-cols-4 + m-1.
export const Portraits = ({
	data,
	width,
	height,
	margin,
	assetBase,
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
					assetBase={assetBase}
				/>
			))}
		</div>
	);
};
