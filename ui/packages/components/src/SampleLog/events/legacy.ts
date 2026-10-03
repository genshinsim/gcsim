import type { LogDetails, Sample } from "@gcsim/types";
import type { EventFields, EventOf, SimEvent } from "./types";

type Logs = LogDetails["logs"];

export function fromLegacySample(sample: Sample): SimEvent[] {
	const lines = parseLines(sample.logs);
	if (lines.length === 0) {
		return [];
	}
	const team = sample.character_details?.map((c) => c.name) ?? [];
	return [
		...onFieldStints(lines, team, sample.initial_character),
		...lines.map(toEvent),
	].sort((a, b) => a.frame - b.frame);
}

/** older samples store the log as a JSON string */
function parseLines(logs: unknown): LogDetails[] {
	let parsed = logs;
	if (typeof logs === "string") {
		try {
			parsed = JSON.parse(logs);
		} catch {
			return [];
		}
	}
	return Array.isArray(parsed) ? parsed : [];
}

function toEvent(line: LogDetails): SimEvent {
	const logs = ordered(line);
	const { ordering: _, ...raw } = { ...line, logs };
	const e: SimEvent = {
		type: line.event,
		frame: line.frame,
		characterIndex: line.char_index,
		message: message(line, logs),
		raw,
	};
	return { ...e, ...endOfLine(line), ...fields(line, logs) };
}

function endOfLine(line: LogDetails): Pick<SimEvent, "end"> {
	if (line.ended > line.frame) {
		return { end: line.ended };
	}
	const neverEnds = line.ended < line.frame;
	return neverEnds ? { end: Number.POSITIVE_INFINITY } : {};
}

function isSwap(line: LogDetails): boolean {
	return (
		line.event === "action" &&
		line.msg.startsWith("executed") &&
		line.logs?.action === "swap"
	);
}

const damageIn = (logs: Logs) => Number(logs.damage) || 0;

function ordered(line: LogDetails): Logs {
	const ordering = line.ordering;
	if (ordering == null) {
		return line.logs ?? {};
	}
	const keys = Object.keys(line.logs ?? {}).sort(
		(a, b) => (ordering[a] ?? 0) - (ordering[b] ?? 0),
	);
	return Object.fromEntries(keys.map((k) => [k, line.logs[k]]));
}

type LineFields =
	| EventFields["action"]
	| EventFields["damage"]
	| EventFields["status"]
	| Record<string, never>;

function fields(line: LogDetails, logs: Logs): LineFields {
	switch (line.event) {
		case "action":
			return { action: typeof logs.action === "string" ? logs.action : "" };
		case "damage":
			return { damage: damageIn(logs), attack: line.msg, mods: modsIn(logs) };
		case "status":
			return { key: typeof logs.key === "string" ? logs.key : "" };
		default:
			return {};
	}
}

function onFieldStints(
	lines: LogDetails[],
	team: string[],
	initial: string | undefined,
): SimEvent[] {
	const last = lines.reduce((m, l) => Math.max(m, l.frame), 0);
	const stints: SimEvent[] = [];
	const addStint = (char: number, start: number, end: number) => {
		if (char < 0 || char >= team.length) {
			return;
		}
		const stint: EventOf<"stint"> = {
			type: "stint",
			frame: start,
			end,
			characterIndex: char,
			message: team[char],
			raw: null,
		};
		stints.push(stint);
	};
	let char = team.indexOf(initial ?? "");
	let start = 0;
	for (const l of lines) {
		if (isSwap(l) && l.char_index !== char) {
			addStint(char, start, l.frame);
			char = l.char_index;
			start = l.frame;
		}
	}
	addStint(char, start, last);
	return stints;
}

/** types whose message reads fine as logged */
const PLAIN = new Set([
	"cooldown",
	"hitlag",
	"enemy",
	"user",
	"calc",
	"character",
	"snapshot",
	"pre_damage_mods",
	"heal",
	"hurt",
	"shield",
	"icd",
	"construct",
	"player",
	"weapon",
	"artifact",
	"warning",
	"debug",
	"sim",
]);

function message(line: LogDetails, d: Logs): string {
	switch (line.event) {
		case "damage":
			return damageMessage(line.msg, d);
		case "action": {
			const msg = isSwap(line) ? `${line.msg} to ${d.target}` : line.msg;
			return msg.replace("executed ", "");
		}
		case "element":
			return elementMessage(line.msg, d);
		case "energy":
			return energyMessage(line.msg, d);
		case "status":
			return `${d.key} ${line.msg}`;
		default:
			return PLAIN.has(line.event) ? line.msg : `${line.event}: ${line.msg}`;
	}
}

const modsIn = (d: Logs) =>
	[d.amp, d.cata, d.crit ? "crit" : ""]
		.filter((x) => typeof x === "string" && x !== "")
		.join(" ");

function damageMessage(msg: string, d: Logs): string {
	const dmg = Math.round(damageIn(d)).toLocaleString("en-US");
	return withMods(`${msg} [${dmg}]`, modsIn(d));
}

export const withMods = (text: string, mods: string) =>
	mods === "" ? text : `${text} (${mods})`;

const auras = (xs: string[]) =>
	xs.map((x) => x.replace(/: (.+)/, " ($1)")).join(" ");

function elementMessage(msg: string, d: Logs): string {
	switch (msg) {
		case "expired":
			return `${d.old_ele} expired`;
		case "application": {
			const before = d.existing ? ` to [${auras(d.existing)}]` : " [no aura]";
			const after = d.after ? ` ➜ [${auras(d.after)}]` : " ➜ [no aura]";
			return `${d.applied_ele} applied${before}${after}`;
		}
		case "refreshed":
			return `${d.ele} refreshed`;
		default:
			return msg;
	}
}

function energyMessage(msg: string, d: Logs): string {
	let out = msg;
	const next = Math.floor(d.post_recovery);
	if (msg.includes("particle")) {
		out = `${msg} from ${d.source}, next: ${next}`;
	}
	if (msg.includes("adding energy")) {
		const amt = d["rec'd"];
		out = `adding ${typeof amt === "number" ? amt.toFixed(2) : amt} energy from ${d.source}, next: ${next}`;
	}
	if (d.max_energy && d.post_recovery === d.max_energy) {
		out += " (max)";
	}
	return out;
}
