import { Badge } from "@gcsim/primitives";
import type { model } from "@gcsim/types";
import { DataColorsConst } from "../../common/gcsim";
import { Avatar } from "../Avatar/Avatar";
import ArtifactsIcon from "./ArtifactsIcon";

const WeaponImage = ({ weapon }: { weapon: model.Weapon }) => {
	return (
		<div className="absolute bottom-[-4px] w-[62] right-[-12px] opacity-85">
			<svg
				key={weapon.name}
				width="62"
				height="65"
				aria-hidden="true"
				onError={(e: React.SyntheticEvent<SVGSVGElement, Event>) => {
					(e.target as SVGAElement).href.baseVal =
						"/api/assets/misc/default.png";
				}}
			>
				<filter id="outlinew">
					<feMorphology
						in="SourceAlpha"
						result="expanded"
						operator="dilate"
						radius="1"
					/>
					<feFlood floodColor="white" />
					<feComposite in2="expanded" operator="in" />
					<feComposite in="SourceGraphic" />
				</filter>
				<filter id="outlineb">
					<feMorphology
						in="SourceAlpha"
						result="expanded"
						operator="dilate"
						radius="1.5"
					/>
					<feFlood floodColor="black" />
					<feComposite in2="expanded" operator="in" />
					<feGaussianBlur stdDeviation="1" />
					<feComposite in="SourceGraphic" />
				</filter>
				<image
					filter="url(#outlinew) url(#outlineb)"
					href={`/api/assets/weapons/${weapon.name}.png`}
					height="55"
					width="55"
					x="0"
					y="3"
				/>
			</svg>
		</div>
	);
};

type CharacterTileProps = {
	char: model.Character | null;
	i: number;
	invalid: boolean;
	onImageLoaded: () => void;
	hideDetails?: boolean;

	//optional classes
	className?: string;
};

export const CharacterTile = ({
	char,
	i,
	invalid,
	onImageLoaded,
	hideDetails = false,
	className = "",
}: CharacterTileProps) => {
	//display an empty card here
	if (char === null) {
		return (
			<div
				className={
					"flex flex-col bg-g-surface-2 border border-g-line rounded-g-sm" +
					(className === "" ? "" : " " + className)
				}
			>
				<div className="flex justify-center">
					<img
						src={"/api/assets/misc/nahida.png"}
						alt=""
						className=" object-contain opacity-50 h-24"
						onLoad={onImageLoaded}
					/>
				</div>
			</div>
		);
	}
	const sets: string[] = [];
	let half = false;

	if (!hideDetails && char.sets && char.sets !== null) {
		for (const [key] of Object.entries(char.sets)) {
			sets.push(key);
		}
		if (sets.length === 1 && char.sets[sets[0]] === 2) {
			half = true;
		}
	}

	return (
		<div
			className={
				"flex flex-col bg-g-surface-2 border border-g-line rounded-g-sm" +
				(className === "" ? "" : " " + className)
			}
		>
			<Avatar
				name={char.name ?? ""}
				element={char.element ?? ""}
				overlay
				onImageLoaded={onImageLoaded}
				className="w-full pt-2 z-0"
			>
				{!hideDetails && char.weapon ? (
					<WeaponImage weapon={char.weapon} />
				) : null}
				<div className="absolute bottom-0 left-0 opacity-85">
					<svg
						width={35}
						height={35}
						aria-hidden="true"
						onError={(e: React.SyntheticEvent<SVGSVGElement, Event>) => {
							(e.target as SVGAElement).href.baseVal =
								"/api/assets/misc/default.png";
						}}
					>
						{sets.length > 0 ? <ArtifactsIcon sets={sets} half={half} /> : null}
					</svg>
				</div>

				{!hideDetails ? (
					<>
						<div
							className={
								"absolute left-[-1px] top-[-1px] flex flex-col gap-0 px-1 py-0 rounded-none " +
								"font-bold font-g-mono text-g-xs rounded-tl-g-sm rounded-br-g-lg " +
								"bg-g-surface-3 opacity-85"
							}
						>
							<div className="flex flex-row gap-1 min-h-fit">
								<span className="text-geo">{`C${char.cons ?? 0}`}</span>
								{!hideDetails && char.weapon ? (
									<span className="text-electro">{`R${
										char.weapon.refine ?? 0
									}`}</span>
								) : null}
							</div>
						</div>
						<div
							className={
								"absolute right-[-1px] top-[-1px] flex flex-col gap-0 px-1 py-0 rounded-none " +
								"font-g-mono text-g-xs rounded-tr-g-sm rounded-bl-g-lg " +
								"bg-g-surface-3 opacity-85"
							}
						>
							<div className="flex flex-row gap-1 min-h-fit items-center">
								<span className="text-g-xs text-g-ink-mute">lvl</span>
								<span
									className={`font-bold`}
									style={{ color: DataColorsConst.qualitative5(i) }}
								>
									{char.level}
								</span>
							</div>
						</div>{" "}
					</>
				) : null}

				{invalid && (
					<div className="absolute left-0 top-1/3 w-full">
						<Badge className="flex flex-row items-center justify-center gap-2">
							<span className="font-g-mono select-none text-g-danger font-bold text-g-xs uppercase">
								WIP
							</span>
						</Badge>
					</div>
				)}
			</Avatar>
		</div>
	);
};
