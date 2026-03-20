import latestCharsData from "./latest-chars.json";

/**
 * Map of version strings to arrays of character key strings added in that version.
 */
export type LatestCharsMap = Record<string, string[]>;

/**
 * Characters added in each game version, keyed by version string (e.g. "v2.38").
 */
export const latestChars: LatestCharsMap = latestCharsData as LatestCharsMap;
