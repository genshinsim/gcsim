import type { model } from "@gcsim/types";
import { GRAY_600, GRAY_700, levelColor, portraitBackground } from "./colors";

type PortraitProps = {
	char: model.Character | null;
	i: number;
	invalid: boolean;
	width: number;
	height: number;
	assetBase: string;
};

// A single character portrait: avatar, weapon, artifact set(s) and level/cons
// badges. All imagery is remote URL <img> behind the injected asset base.
const Portrait = ({
	char,
	i,
	invalid,
	width,
	height,
	assetBase,
}: PortraitProps) => {
	const base = {
		display: "flex" as const,
		width,
		height,
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
					width={height * 0.9}
					height={height * 0.9}
					alt=""
					style={{ objectFit: "contain", opacity: 0.5 }}
				/>
			</div>
		);
	}

	const sets = char.sets ? Object.keys(char.sets) : [];
	// A lone 2-piece set renders as a half-width flower (matches AvatarPortrait).
	const isHalfWidthSet = sets.length === 1 && char.sets?.[sets[0]] === 2;

	const avatarSize = height * 0.92;
	const badgeStyle = {
		display: "flex" as const,
		flexDirection: "row" as const,
		gap: 3,
		padding: "0 3px",
		backgroundColor: GRAY_700,
		fontFamily: "Inter",
		fontSize: 10,
		fontWeight: 700,
	};

	return (
		<div
			style={{
				...base,
				alignItems: "center",
				justifyContent: "center",
				backgroundImage: portraitBackground(char.element ?? ""),
			}}
		>
			{/* avatar */}
			<img
				src={`${assetBase}/avatar/${char.name}.png`}
				width={avatarSize}
				height={avatarSize}
				alt={char.name ?? ""}
				style={{ objectFit: "contain" }}
			/>

			{/* weapon */}
			{char.weapon?.name ? (
				<img
					src={`${assetBase}/weapons/${char.weapon.name}.png`}
					width={height * 0.5}
					height={height * 0.5}
					alt=""
					style={{
						position: "absolute",
						bottom: -2,
						right: -6,
						objectFit: "contain",
						opacity: 0.85,
					}}
				/>
			) : null}

			{/* artifact set(s) */}
			{sets.length > 0 ? (
				<div
					style={{
						position: "absolute",
						bottom: 1,
						left: 1,
						display: "flex",
						flexDirection: "row",
						opacity: 0.9,
					}}
				>
					<img
						src={`${assetBase}/artifacts/${sets[0]}_flower.png`}
						width={sets.length > 1 || isHalfWidthSet ? 16 : 32}
						height={32}
						alt=""
						style={{ objectFit: "cover" }}
					/>
					{sets.length > 1 ? (
						<img
							src={`${assetBase}/artifacts/${sets[1]}_flower.png`}
							width={16}
							height={32}
							alt=""
							style={{ objectFit: "cover" }}
						/>
					) : null}
				</div>
			) : null}

			{/* cons / refine */}
			<div
				style={{
					...badgeStyle,
					position: "absolute",
					left: 0,
					top: 0,
					borderBottomRightRadius: 6,
				}}
			>
				<span style={{ color: "#facc15" }}>{`C${char.cons ?? 0}`}</span>
				{char.weapon ? (
					<span
						style={{ color: "#c084fc" }}
					>{`R${char.weapon.refine ?? 0}`}</span>
				) : null}
			</div>

			{/* level */}
			<div
				style={{
					...badgeStyle,
					position: "absolute",
					right: 0,
					top: 0,
					borderBottomLeftRadius: 6,
					alignItems: "center",
				}}
			>
				<span style={{ color: "#9ca3af", fontSize: 9 }}>lvl</span>
				<span style={{ color: levelColor(i) }}>{char.level}</span>
			</div>

			{invalid ? (
				<div
					style={{
						position: "absolute",
						top: "33%",
						left: 0,
						width,
						display: "flex",
						justifyContent: "center",
					}}
				>
					<span
						style={{
							color: "#ef4444",
							fontFamily: "Inter",
							fontWeight: 700,
							fontSize: 11,
							backgroundColor: GRAY_700,
							padding: "0 4px",
							borderRadius: 4,
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
	gap: number;
	assetBase: string;
};

// The four-portrait row (flexbox, replacing the live card's CSS grid).
export const Portraits = ({ data, width, height, gap, assetBase }: Props) => {
	const chars = data.character_details ?? [];
	return (
		<div
			style={{
				display: "flex",
				flexDirection: "row",
				gap,
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
					assetBase={assetBase}
				/>
			))}
		</div>
	);
};
