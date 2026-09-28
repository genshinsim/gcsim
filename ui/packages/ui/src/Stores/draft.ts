import {
	createStore,
	type KeyStorage,
	LEGACY_APP_DATA_KEY,
	readJSON,
	type Store,
} from "./externalStore";

const DRAFT_KEY = "gcsim-draft-config";

const charLinesRegEx =
	/^(\w+) (?:char|add) (?:lvl|weapon|set|stats).+$(?:\r\n|\r|\n)?/gm;

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

export interface DraftStore extends Store<string> {
	send(cfg: string, opts: { keepTeam: boolean }): void;
}

function loadDraft(storage: KeyStorage): string {
	return (
		storage.getItem(DRAFT_KEY) ?? readJSON(storage, LEGACY_APP_DATA_KEY)?.cfg ?? ""
	);
}

export function draftStore(storage: KeyStorage): DraftStore {
	const store = createStore(loadDraft(storage), (cfg) =>
		storage.setItem(DRAFT_KEY, cfg),
	);
	return {
		...store,
		send: (cfg, { keepTeam }) =>
			store.set((current) => (keepTeam ? mergeTeam(current, cfg) : cfg)),
	};
}
