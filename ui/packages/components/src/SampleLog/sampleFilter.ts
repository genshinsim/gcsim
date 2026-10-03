import { useCallback, useState } from "react";
import { AllSampleOptions } from "./SampleOptions";

const FILTER_KEY = "gcsim-sample-log-filter";

type FilterStorage = Pick<Storage, "getItem" | "setItem">;

const isFilter = (v: unknown): v is string[] =>
	Array.isArray(v) && v.every((x) => typeof x === "string");

export function loadSampleFilter(storage: FilterStorage): string[] {
	try {
		const saved: unknown = JSON.parse(storage.getItem(FILTER_KEY) ?? "null");
		return isFilter(saved) ? saved : AllSampleOptions;
	} catch {
		return AllSampleOptions;
	}
}

export function saveSampleFilter(storage: FilterStorage, val: string[]) {
	storage.setItem(FILTER_KEY, JSON.stringify(val));
}

export function useSampleFilter() {
	const [filter, setFilterState] = useState(() =>
		loadSampleFilter(localStorage),
	);
	const setFilter = useCallback((val: string[]) => {
		setFilterState(val);
		saveSampleFilter(localStorage, val);
	}, []);
	return [filter, setFilter] as const;
}
