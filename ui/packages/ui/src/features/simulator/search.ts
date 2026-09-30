export type SimulatorSearch = {
	cfg?: string;
};

export function validateSimulatorSearch(
	raw: Record<string, unknown>,
): SimulatorSearch {
	return typeof raw.cfg === "string" && raw.cfg !== "" ? { cfg: raw.cfg } : {};
}
