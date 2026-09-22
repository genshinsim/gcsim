import {
	Button,
	Tooltip,
	TooltipContent,
	TooltipProvider,
	TooltipTrigger,
} from "@gcsim/primitives";
import type { model } from "@gcsim/types";
import {
	ChevronDownIcon,
	ChevronUpIcon,
	SearchIcon,
	XIcon,
	ZoomInIcon,
} from "lucide-react";
import type { JSX } from "react";
import { charCardBG } from "../../lib/helper";
import { Avatar } from "../Avatar/Avatar";
import {
	IconAnemo,
	IconAtk,
	IconCD,
	IconCR,
	IconCryo,
	IconDef,
	IconDendro,
	IconElectro,
	IconEM,
	IconER,
	IconGeo,
	IconHeal,
	IconHP,
	IconHydro,
	IconPhysical,
	IconPyro,
} from "./Icons";
import { WeaponCard } from "./WeaponCard";

export type CharStatBlock = {
	key: string;
	name: string;
	t: string; // "both" | "f" | "%"
	flat: number;
	percent: number;
};

type Props = {
	char: model.Character;
	stats: CharStatBlock[];
	snapshot: CharStatBlock[];
	statsRows: number;

	name: string;
	constellationLabel: string;
	levelLabel: string;
	talentsLabel: string;
	artifactStatsLabel: string;
	totalStatsLabel: string;
	weaponName: string;

	className?: string;
	showDetails?: boolean;
	showSnapshot?: boolean;
	viewerMode?: boolean;
	isSkeleton?: boolean;
	handleDelete?: () => void;
	handleToggleDetail?: () => void;
	handleToggleSnapshot?: () => void;
};

function statKeyToIcon(key: string): JSX.Element {
	switch (key) {
		case "hp":
			return <IconHP />;
		case "atk":
			return <IconAtk />;
		case "def":
			return <IconDef />;
		case "er":
			return <IconER />;
		case "em":
			return <IconEM />;
		case "cr":
			return <IconCR />;
		case "cd":
			return <IconCD />;
		case "electro":
			return <IconElectro />;
		case "pyro":
			return <IconPyro />;
		case "cryo":
			return <IconCryo />;
		case "hydro":
			return <IconHydro />;
		case "geo":
			return <IconGeo />;
		case "anemo":
			return <IconAnemo />;
		case "phys":
			return <IconPhysical />;
		case "heal":
			return <IconHeal />;
		case "dendro":
			return <IconDendro />;
		default:
			return <span />;
	}
}

export function CharacterCard({
	char,
	stats,
	snapshot,
	statsRows,
	name,
	constellationLabel,
	levelLabel,
	talentsLabel,
	artifactStatsLabel,
	totalStatsLabel,
	weaponName,
	showDetails = true,
	showSnapshot = true,
	viewerMode = false,
	isSkeleton,
	handleDelete,
	handleToggleDetail,
	handleToggleSnapshot,
	className = "",
}: Props) {
	const arts: JSX.Element[] = [];

	for (const key in char.sets) {
		arts.push(
			<div className="w-8 flex flex-col rounded-g-md" key={key}>
				<Tooltip>
					<TooltipTrigger asChild>
						<img
							src={`/api/assets/artifacts/${key}_flower.png`}
							alt={key}
							className="w-full h-8"
							onError={(e) =>
								((e.target as HTMLImageElement).src =
									"/api/assets/misc/default.png")
							}
						/>
					</TooltipTrigger>
					<TooltipContent>{key}</TooltipContent>
				</Tooltip>

				<span className="text-center text-g-xs">{char.sets?.[key]}</span>
			</div>,
		);
	}

	let count = 0;
	const rows: JSX.Element[] = [];

	let statsHeader = artifactStatsLabel;
	if (showSnapshot && viewerMode) {
		stats = snapshot;
		statsHeader = totalStatsLabel;
	}
	stats.forEach((s) => {
		const val: JSX.Element[] = [];
		if (s.flat === 0 && s.percent === 0) {
			return;
		}

		count++;

		switch (s.t) {
			case "both":
				val.push(
					<td key={"flat-" + s.key} className="text-right text-g-xs">
						{s.flat.toFixed(0)}
					</td>,
				);
				val.push(
					<td key={"per-" + s.key} className="text-right text-g-xs">
						{(s.percent * 100).toFixed(2) + "%"}
					</td>,
				);
				break;
			case "f":
				val.push(
					<td key={"flat-" + s.key} className="text-right text-g-xs">
						{s.flat.toFixed(0)}
					</td>,
				);
				val.push(<td key={"per-" + s.key}></td>);
				break;
			case "%":
				val.push(<td key={"flat-" + s.key}></td>);
				val.push(
					<td key={"per-" + s.key} className="text-right text-g-xs">
						{(s.percent * 100).toFixed(2) + "%"}
					</td>,
				);
		}

		rows.push(
			<tr key={count}>
				<td className="flex flex-row gap-0.5 place-items-center">
					<div className="w-4 mr-1 fill-g-ink">{statKeyToIcon(s.key)}</div>
					<span className="text-g-xs sm:text-g-sm">{s.name}</span>
				</td>
				{val}
			</tr>,
		);
	});

	for (; count < statsRows; count++) {
		rows.push(
			<tr key={count + 1}>
				<td>
					<br />
				</td>
				<td></td>
				<td></td>
			</tr>,
		);
	}

	const skeleton = isSkeleton
		? "animate-pulse rounded bg-g-surface-2 text-transparent"
		: "";

	return (
		<div className={className}>
			<div className="min-h-24 bg-g-surface text-g-ink shadow text-g-sm flex flex-col justify-center gap-2 border border-g-line">
				<div
					className={
						"character-parent flex flex-row pt-4 pl-4 pr-2 relative z-0 " +
						charCardBG(char.element ?? "")
					}
				>
					<div className="flex flex-row gap-1 absolute top-1 right-1">
						<div className="flex flex-col gap-1">
							<Button
								variant="secondary"
								size="icon-xs"
								onClick={handleToggleDetail}
							>
								{showDetails ? <ChevronUpIcon /> : <ChevronDownIcon />}
							</Button>
							{showDetails && viewerMode ? (
								<Button
									variant="secondary"
									size="icon-xs"
									onClick={handleToggleSnapshot}
								>
									{showSnapshot ? <SearchIcon /> : <ZoomInIcon />}
								</Button>
							) : null}
						</div>
						{viewerMode ? null : (
							<Button
								variant="destructive"
								size="icon-xs"
								onClick={handleDelete}
							>
								<XIcon />
							</Button>
						)}
					</div>
					<div
						className="character-header absolute inset-0 -z-10 !bg-cover !bg-center mix-blend-luminosity opacity-75"
						style={{ background: "url(/api/assets/misc/overlay.jpg)" }}
					></div>
					<div
						className={
							"character-name font-medium m-4 capitalize absolute top-0 left-0 " +
							skeleton
						}
					>
						{constellationLabel} {name}
					</div>
					<div className="w-1/2 text-g-sm">
						<div className={"pl-1 pr-1 mt-6 " + skeleton}>
							<div>
								{levelLabel} {char.level}/{char.max_level}
							</div>
							<div>
								{talentsLabel} {char.talents?.attack}/{char.talents?.skill}/
								{char.talents?.burst}
							</div>
							<TooltipProvider>
								<div className="mt-1 mr-2 grid grid-cols-5">{arts}</div>
							</TooltipProvider>
						</div>
					</div>
					<div className="w-1/2 h-32">
						{isSkeleton ? null : (
							<Avatar
								name={char.name ?? ""}
								background={false}
								imageWrapClassName=""
								imageClassName="ml-auto h-32"
							/>
						)}
					</div>
				</div>

				{char.weapon ? (
					<WeaponCard
						weapon={char.weapon}
						name={weaponName}
						isSkeleton={isSkeleton}
					/>
				) : null}

				{showDetails ? (
					<div className="flex flex-col gap-2 mx-2 p-2 bg-g-surface-2 border-g-line border">
						<span className="font-bold">{statsHeader}</span>
						<div className="px-2">
							<table className="w-full">
								<tbody>{rows}</tbody>
							</table>
						</div>
					</div>
				) : null}

				<div className="" />
			</div>
		</div>
	);
}
