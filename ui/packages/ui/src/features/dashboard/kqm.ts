import { model } from "@gcsim/types";

export const KQM_DB_URL = "https://db.kqm.gg";

export type KqmEntry = {
	_id?: string;
	description?: string;
	submitter?: string;
	summary?: {
		team?: ({ name?: string } | null)[];
		mean_dps_per_target?: number;
		mode?: model.SimMode;
	};
};

type JsonObject = Record<string, unknown>;

const isObject = (v: unknown): v is JsonObject =>
	typeof v === "object" && v !== null && !Array.isArray(v);

const str = (v: unknown) => (typeof v === "string" ? v : undefined);

function parseSummary(s: JsonObject): NonNullable<KqmEntry["summary"]> {
	return {
		team: Array.isArray(s.team)
			? s.team.map((c) => (isObject(c) ? { name: str(c.name) } : null))
			: undefined,
		mean_dps_per_target:
			typeof s.mean_dps_per_target === "number"
				? s.mean_dps_per_target
				: undefined,
		mode: s.mode == null ? undefined : model.simModeFromJSON(s.mode),
	};
}

function parseEntry(e: JsonObject): KqmEntry {
	return {
		_id: str(e._id) ?? str(e.id),
		description: str(e.description),
		submitter: str(e.submitter),
		summary: isObject(e.summary) ? parseSummary(e.summary) : undefined,
	};
}

export function parseKqmEntries(json: unknown): KqmEntry[] {
	if (!isObject(json) || !Array.isArray(json.data)) return [];
	return json.data.filter(isObject).map(parseEntry);
}
