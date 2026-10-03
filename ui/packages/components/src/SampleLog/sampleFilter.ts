import React from "react";
import { DefaultSampleOptions } from "./SampleOptions";

const FILTER_KEY = "gcsim-sample-settings";

type FilterStorage = Pick<Storage, "getItem" | "setItem">;

export function loadSampleFilter(storage: FilterStorage): string[] {
	try {
		const raw = storage.getItem(FILTER_KEY);
		return (raw && JSON.parse(raw)) || DefaultSampleOptions;
	} catch {
		return DefaultSampleOptions;
	}
}

export function saveSampleFilter(storage: FilterStorage, val: string[]) {
	storage.setItem(FILTER_KEY, JSON.stringify(val));
}

export function useSampleFilter() {
	const [filter, setFilter] = React.useState(() =>
		loadSampleFilter(localStorage),
	);
	React.useEffect(() => saveSampleFilter(localStorage, filter), [filter]);
	return [filter, setFilter] as const;
}
