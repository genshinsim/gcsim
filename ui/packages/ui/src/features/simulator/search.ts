export type SimulatorSearch = {
	// base64 of a gzipped JSON SharedConfig; see sharedConfig.ts
	cfg?: string;
};

export function validateSimulatorSearch(
	raw: Record<string, unknown>,
): SimulatorSearch {
	return typeof raw.cfg === "string" && raw.cfg !== "" ? { cfg: raw.cfg } : {};
}
