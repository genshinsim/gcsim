import type { LogDetails, Sample } from "@gcsim/types";
import type { EventFields, SimEvent } from "./types";

type Logs = LogDetails["logs"];

export function fromLegacySample(sample: Sample): SimEvent[] {
	const lines = parseLines(sample.logs as LogDetails[] | string | undefined);
	if (lines.length === 0) {
		return [];
	}
	const team = sample.character_details?.map((c) => c.name) ?? [];
	return [
		...onFieldStints(lines, team, sample.initial_character),
		...lines.map(toEvent),
	].sort((a, b) => a.frame - b.frame);
}

function parseLines(logs: LogDetails[] | string | undefined): LogDetails[] {
	if (logs == null) {
		return [];
	}
	if (typeof logs !== "string") {
		return logs;
	}
	try {
		return JSON.parse(logs);
	} catch {
		return [];
	}
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
	return { ...e, ...endOfLine(line), ...fields(line.event, logs) };
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

const damageOf = (logs: Logs) => Number(logs.damage) || 0;

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

function fields(type: string, logs: Logs): object {
	switch (type) {
		case "action":
			return {
				action: typeof logs.action === "string" ? logs.action : "",
			} satisfies EventFields["action"];
		case "damage":
			return {
				damage: damageOf(logs),
			} satisfies EventFields["damage"];
		case "status":
			return {
				key: typeof logs.key === "string" ? logs.key : "",
			} satisfies EventFields["status"];
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
		if (char < 0) {
			return;
		}
		const stint: SimEvent & EventFields["stint"] = {
			type: "stint",
			frame: start,
			end,
			characterIndex: char,
			message: team[char] ?? "",
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
		case "cooldown":
		case "hitlag":
		case "enemy":
		case "user":
		case "calc":
		case "character":
		case "snapshot":
		case "pre_damage_mods":
		case "heal":
		case "hurt":
		case "shield":
		case "icd":
		case "construct":
		case "player":
		case "weapon":
		case "artifact":
		case "warning":
		case "debug":
		case "sim":
			return line.msg;
		default:
			return `${line.event}: ${line.msg}`;
	}
}

function damageMessage(msg: string, d: Logs): string {
	const dmg = Math.round(damageOf(d)).toLocaleString("en-US");
	const extra = [d.amp, d.cata, d.crit ? "crit" : ""]
		.filter((x) => typeof x === "string" && x !== "")
		.join(" ");
	return `${msg} [${dmg}]${extra === "" ? "" : ` (${extra})`}`;
}

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
