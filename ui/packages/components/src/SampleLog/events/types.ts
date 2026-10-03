export interface SimEvent {
	type: string;
	frame: number;
	/** absent for an instant event; Infinity when it never ends */
	end?: number;
	characterIndex: number;
	message: string;
	raw: unknown;
}

export interface EventFields {
	action: { action: string };
	damage: { damage: number };
	status: { key: string };
	stint: { end: number };
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
