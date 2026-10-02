import type { SavedRunStore } from "@gcsim/components";
import type { KeyStorage } from "./storage";

const LEGACY_RESULT_KEY = "redux-local-results";
const LEGACY_HASH_KEY = "redux-local-results-hash";
const RAW_KEY = "local-results-raw";
const HASH_KEY = "local-results-raw-hash";

export function lastRunStore(storage: KeyStorage): SavedRunStore {
	return {
		load: () => {
			const raw = storage.getItem(RAW_KEY);
			const legacy = storage.getItem(LEGACY_RESULT_KEY);
			try {
				if (raw) {
					return {
						result: JSON.parse(raw),
						raw,
						hash: storage.getItem(HASH_KEY) || null,
					};
				}
				if (legacy) {
					return { result: JSON.parse(legacy), raw: null, hash: null };
				}
			} catch {}
			return null;
		},
		save: ({ result, raw, hash }) => {
			if (raw == null) {
				storage.setItem(LEGACY_RESULT_KEY, JSON.stringify(result));
				storage.setItem(RAW_KEY, "");
			} else {
				storage.setItem(RAW_KEY, raw);
				storage.setItem(LEGACY_RESULT_KEY, "");
			}
			storage.setItem(HASH_KEY, hash ?? "");
			storage.setItem(LEGACY_HASH_KEY, "");
		},
	};
}
