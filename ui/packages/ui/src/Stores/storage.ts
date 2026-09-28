export type KeyStorage = Pick<Storage, "getItem" | "setItem">;

export const LEGACY_APP_DATA_KEY = "redux-app-data";

export function readJSON(storage: KeyStorage, key: string) {
	try {
		const raw = storage.getItem(key);
		return raw ? JSON.parse(raw) : null;
	} catch {
		return null;
	}
}
