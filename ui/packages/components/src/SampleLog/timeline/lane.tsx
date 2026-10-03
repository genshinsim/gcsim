import { dynamicKey } from "@gcsim/localization";
import { cn } from "@gcsim/primitives";
import { Swords } from "lucide-react";
import type React from "react";
import { useCallback } from "react";
import { useTranslation } from "react-i18next";
import { Avatar } from "../../Cards/Avatar/Avatar";
import { isSimLane, type Lane } from "./model";

export const laneColorVar = (lane: Lane) =>
	lane.element == null ? "--g-text-mute" : `--g-${lane.element}`;

export const laneColor = (lane: Lane) => `var(${laneColorVar(lane)})`;

export const laneTint = (lane: Lane, percent: number) =>
	`color-mix(in srgb, ${laneColor(lane)} ${percent}%, transparent)`;

export function useLaneName() {
	const { t } = useTranslation();
	return useCallback(
		(lane: Lane) =>
			isSimLane(lane)
				? t("sample.sim_lane")
				: t(dynamicKey(`game:character_names.${lane.key}`)),
		[t],
	);
}

/** the character's portrait, or crossed swords for the sim lane; `className` and `style` size the portrait only */
export function LaneIcon({
	lane,
	name,
	className,
	style,
}: {
	lane: Lane;
	name: string;
	className?: string;
	style?: React.CSSProperties;
}) {
	if (isSimLane(lane)) {
		return <Swords className="size-4 text-g-ink-mute" aria-label={name} />;
	}
	return (
		<div className={cn("shrink-0", className)} style={style}>
			<Avatar
				name={lane.key}
				element={lane.element ?? ""}
				className="size-full overflow-hidden rounded-g-md"
				imageClassName="size-full object-contain"
				imageWrapClassName="flex size-full"
			/>
		</div>
	);
}
