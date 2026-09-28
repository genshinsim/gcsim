import { type KeyStorage, LEGACY_APP_DATA_KEY, readJSON } from "./storage";

const DRAFT_KEY = "gcsim-draft-config";

const charLinesRegEx =
	/^(\w+) (?:char|add) (?:lvl|weapon|set|stats).+$(?:\r\n|\r|\n)?/gm;

export type SendOptions = { keepTeam: boolean };

export function mergeTeam(current: string, incoming: string): string {
	let team = "";
	let lastChar = "";
	for (const match of current.matchAll(charLinesRegEx)) {
		if (match[1] !== lastChar) {
			team += "\n";
			lastChar = match[1];
		}
		team += match[0];
	}
	const next = `${team}\n${incoming.replace(charLinesRegEx, "")}`;
	return next.replace(/(\r\n|\r|\n){2,}/g, "$1\n");
}

export function sendDraft(
	current: string,
	incoming: string,
	{ keepTeam }: SendOptions,
): string {
	return keepTeam ? mergeTeam(current, incoming) : incoming;
}

export function loadDraft(storage: KeyStorage): string {
	return (
		storage.getItem(DRAFT_KEY) ??
		readJSON(storage, LEGACY_APP_DATA_KEY)?.cfg ??
		""
	);
}

export function saveDraft(storage: KeyStorage, cfg: string) {
	storage.setItem(DRAFT_KEY, cfg);
}
