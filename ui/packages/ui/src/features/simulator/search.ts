export type SimulatorSearch = {
	// A full config string to inject into the editor on load, e.g. from an
	// external tool that generated a config and wants an "open in gcsim"
	// link. No custom encoding is needed for newlines/quotes/etc: the router
	// skips JSON-wrapping plain string values, so a multi-line cfg round-trips
	// through ordinary URL percent-encoding untouched.
	cfg?: string;
};

export function validateSimulatorSearch(
	raw: Record<string, unknown>,
): SimulatorSearch {
	return typeof raw.cfg === "string" && raw.cfg !== "" ? { cfg: raw.cfg } : {};
}
