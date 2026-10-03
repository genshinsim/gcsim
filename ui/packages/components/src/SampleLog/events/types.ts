/**
 * One thing that happened in a sample, independent of how gcsim logs it.
 * Views read these; only an adapter (see legacy.ts) reads the raw log.
 */
export interface SimEvent {
	type: string;
	frame: number;
	/** absent = instant; null = lasts past the end of the sample */
	end?: number | null;
	/** -1 = sim/global, 0..N-1 = character */
	characterIndex: number;
	message: string;
	/** source data, only for the details view */
	raw: unknown;
}

/**
 * Fields carried by events of a given type. Add an entry when a view needs a
 * field to draw with; an adapter emitting that type must fill it.
 */
export interface EventFields {
	action: { action: string };
	damage: { damage: number };
	status: { key: string };
	/** a character's stint on field */
	field: { end: number };
}

export type EventOf<K extends keyof EventFields> = SimEvent & {
	type: K;
} & EventFields[K];

export function isEvent<K extends keyof EventFields>(
	e: SimEvent,
	type: K,
): e is EventOf<K> {
	return e.type === type;
}
