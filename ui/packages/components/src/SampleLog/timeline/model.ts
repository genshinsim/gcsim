import type { Sample } from "@gcsim/types";
import { fromLegacySample } from "../events/legacy";
import { isEvent, type SimEvent } from "../events/types";

export type TeamMember = { key: string; element: string | null };

export type Lane = TeamMember & {
	index: number;
	/** "sim" for the sim/enemy lane, else the character key */
	key: string;
};

export type Chip = {
	id: number;
	lane: number;
	/** the column it sits in: the event's frame, or its end for an expiry chip */
	frame: number;
	event: SimEvent;
	/** marks where the event ended rather than where it started */
	expired: boolean;
	dmg: number;
};

export type Stint = { lane: number; start: number; end: number };

export type TimelineModel = {
	lanes: Lane[];
	chips: Chip[];
	onField: Stint[];
	maxFrame: number;
	counts: Map<string, number>;
};

export function chipText(c: Chip): string {
	if (!c.expired) {
		return c.event.message;
	}
	return `${spanLabel(c.event)} expired`;
}

/** short name for the thing that lasted, e.g. a status key */
export function spanLabel(e: SimEvent): string {
	return isEvent(e, "status") && e.key !== "" ? e.key : e.message;
}

export function modelFromSample(sample: Sample): TimelineModel {
	return modelFromEvents(
		fromLegacySample(sample),
		sample.character_details?.map((c) => ({
			key: c.name,
			element: c.element ?? null,
		})) ?? [],
	);
}

export function modelFromEvents(
	all: SimEvent[],
	team: TeamMember[],
): TimelineModel {
	const lanes: Lane[] = [
		{ index: 0, key: "sim", element: null },
		...team.map((c, i) => ({ ...c, index: i + 1 })),
	];
	const laneOf = (e: SimEvent) =>
		e.characterIndex >= 0 && e.characterIndex < team.length
			? e.characterIndex + 1
			: 0;
	let maxFrame = 0;
	for (const e of all) {
		maxFrame = Math.max(maxFrame, e.frame);
	}
	const expiresInSample = (e: SimEvent): e is SimEvent & { end: number } =>
		isEvent(e, "status") && e.key !== "" && e.end != null && e.end <= maxFrame;

	const chips: Chip[] = [];
	const onField: Stint[] = [];
	const counts = new Map<string, number>();
	const addChip = (e: SimEvent, frame: number, expired: boolean) =>
		chips.push({
			id: chips.length,
			lane: laneOf(e),
			frame,
			event: e,
			expired,
			dmg: !expired && isEvent(e, "damage") ? e.damage : 0,
		});
	for (const e of all) {
		if (isEvent(e, "field")) {
			onField.push({ lane: laneOf(e), start: e.frame, end: e.end });
			continue;
		}
		addChip(e, e.frame, false);
		counts.set(e.type, (counts.get(e.type) ?? 0) + 1);
		if (expiresInSample(e)) {
			addChip(e, e.end, true);
		}
	}
	chips.sort((a, b) => a.frame - b.frame || a.id - b.id);

	return { lanes, chips, onField, maxFrame, counts };
}

/** index of first element with key >= value */
export function lowerBound<T>(arr: T[], value: number, key: (t: T) => number) {
	let lo = 0;
	let hi = arr.length;
	while (lo < hi) {
		const mid = (lo + hi) >> 1;
		if (key(arr[mid]) < value) {
			lo = mid + 1;
		} else {
			hi = mid;
		}
	}
	return lo;
}

export function onFieldLaneAt(onField: Stint[], frame: number): number {
	const i = lowerBound(onField, frame + 1, (s) => s.start) - 1;
	return i >= 0 ? onField[i].lane : -1;
}
