import type { model, ParsedResult } from "@gcsim/types";

export function toParsedTeam(result: ParsedResult): model.Character[] {
	if (!result.characters) {
		return [];
	}
	return result.characters.map((c) => ({
		name: c.base.key,
		level: c.base.level,
		element: c.base.element,
		max_level: c.base.max_level,
		cons: c.base.cons,
		weapon: c.weapon,
		talents: c.talents,
		stats: c.stats,
		snapshot: new Array(c.stats?.length ?? 0).fill(0),
		sets: c.sets,
	}));
}
