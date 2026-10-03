import React from "react";
import { DefaultSampleOptions } from "./SampleOptions";

const FILTER_KEY = "gcsim-sample-settings";

type FilterStorage = Pick<Storage, "getItem" | "setItem">;

const isFilter = (v: unknown): v is string[] =>
	Array.isArray(v) && v.every((x) => typeof x === "string");

export function loadSampleFilter(storage: FilterStorage): string[] {
	try {
		const saved: unknown = JSON.parse(storage.getItem(FILTER_KEY) ?? "null");
		return isFilter(saved) ? saved : DefaultSampleOptions;
	} catch {
		return DefaultSampleOptions;
	}
}

export function saveSampleFilter(storage: FilterStorage, val: string[]) {
	storage.setItem(FILTER_KEY, JSON.stringify(val));
}

export function useSampleFilter() {
	const [filter, setFilterState] = React.useState(() =>
		loadSampleFilter(localStorage),
	);
	const setFilter = React.useCallback((val: string[]) => {
		setFilterState(val);
		saveSampleFilter(localStorage, val);
	}, []);
	return [filter, setFilter] as const;
}
