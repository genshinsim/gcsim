import { useCallback, useState } from "react";

const VIEW_KEY = "gcsim-sample-log-view";

export const CAP_MIN = 3;
export const CAP_MAX = 10;

export type Orientation = "horizontal" | "vertical";

export type SampleView = {
	orientation: Orientation;
	/** one-line chips capped per row; otherwise chips wrap and every one shows */
	slim: boolean;
	/** chips a vertical row shows before "+N more" when slim */
	cap: number;
	/** the vertical layout's overview is open */
	overview: boolean;
};

type ViewStorage = Pick<Storage, "getItem" | "setItem">;

function parse(storage: ViewStorage): Record<string, unknown> {
	try {
		const v: unknown = JSON.parse(storage.getItem(VIEW_KEY) ?? "null");
		return typeof v === "object" && v != null ? { ...v } : {};
	} catch {
		return {};
	}
}

export function loadSampleView(
	storage: ViewStorage,
	overviewByDefault: boolean,
): SampleView {
	const v = parse(storage);
	return {
		orientation: v.orientation === "vertical" ? "vertical" : "horizontal",
		slim: typeof v.slim === "boolean" ? v.slim : true,
		cap:
			typeof v.cap === "number"
				? Math.min(CAP_MAX, Math.max(CAP_MIN, Math.round(v.cap)))
				: 4,
		overview: typeof v.overview === "boolean" ? v.overview : overviewByDefault,
	};
}

export function saveSampleView(storage: ViewStorage, view: SampleView) {
	storage.setItem(VIEW_KEY, JSON.stringify(view));
}

export function useSampleView() {
	const [view, setViewState] = useState(() =>
		loadSampleView(
			localStorage,
			window.matchMedia?.("(min-width: 640px)").matches ?? true,
		),
	);
	const setView = useCallback((v: SampleView) => {
		setViewState(v);
		saveSampleView(localStorage, v);
	}, []);
	return [view, setView] as const;
}
