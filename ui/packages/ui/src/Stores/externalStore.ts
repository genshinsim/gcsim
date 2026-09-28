import { useSyncExternalStore } from "react";

export type KeyStorage = Pick<Storage, "getItem" | "setItem">;

export interface Store<T> {
	get(): T;
	set(next: T | ((prev: T) => T)): void;
	subscribe(listener: () => void): () => void;
}

export function createStore<T>(
	initial: T,
	onChange?: (value: T) => void,
): Store<T> {
	let value = initial;
	const listeners = new Set<() => void>();
	return {
		get: () => value,
		set: (next) => {
			value = next instanceof Function ? next(value) : next;
			onChange?.(value);
			for (const l of listeners) {
				l();
			}
		},
		subscribe: (listener) => {
			listeners.add(listener);
			return () => listeners.delete(listener);
		},
	};
}

export function useStore<T>(store: Store<T>): T {
	return useSyncExternalStore(store.subscribe, store.get);
}

export function readJSON(storage: KeyStorage, key: string) {
	try {
		const raw = storage.getItem(key);
		return raw ? JSON.parse(raw) : null;
	} catch {
		return null;
	}
}

export const legacyAppDataKey = "redux-app-data";
