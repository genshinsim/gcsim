export function asError(err: unknown): string {
	return typeof err === "string" ? err : String(err);
}
