export const VIEWER_TABS = ["results", "config", "sample"] as const;

export type ViewerTab = (typeof VIEWER_TABS)[number];

export type ViewerSearch = {
	tab?: ViewerTab;
	seed?: string;
};

export function validateViewerSearch(
	raw: Record<string, unknown>,
): ViewerSearch {
	const search: ViewerSearch = {};
	if (VIEWER_TABS.includes(raw.tab as ViewerTab)) {
		search.tab = raw.tab as ViewerTab;
	}
	const seed = typeof raw.seed === "number" ? String(raw.seed) : raw.seed;
	if (typeof seed === "string" && /^-?\d+$/.test(seed)) {
		search.seed = seed;
	}
	return search;
}

export function legacyHashSearch(hash: string): ViewerSearch | null {
	const params = new URLSearchParams(hash);
	if (!params.has("tab") && !params.has("sample")) {
		return null;
	}
	return validateViewerSearch({
		tab: params.get("tab"),
		seed: params.get("sample"),
	});
}
