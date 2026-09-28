import type { RunResultStore } from "@gcsim/components";

const RESULT_KEY = "redux-local-results";
const HASH_KEY = "redux-local-results-hash";

type RunStorage = Pick<Storage, "getItem" | "setItem">;

export function lastRunStore(storage: RunStorage): RunResultStore {
	return {
		load: () => {
			const raw = storage.getItem(RESULT_KEY);
			if (!raw) {
				return null;
			}
			try {
				return { result: JSON.parse(raw), hash: storage.getItem(HASH_KEY) };
			} catch {
				return null;
			}
		},
		save: ({ result, hash }) => {
			storage.setItem(RESULT_KEY, JSON.stringify(result));
			storage.setItem(HASH_KEY, hash ?? "");
		},
	};
}
