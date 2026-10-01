import data from "./keys.gen.json";
import artifacts from "./names/artifact.dm.json";
import characters from "./names/character.dm.json";
import weapons from "./names/weapon.dm.json";

export const KEYWORDS: readonly string[] = data.keywords;
export const STATS: readonly string[] = data.stats;
export const ELEMENTS: readonly string[] = data.elements;
export const ACTIONS: readonly string[] = data.actions;

export const CHARACTERS: Readonly<Record<string, string>> = characters;
export const WEAPONS: Readonly<Record<string, string>> = weapons;
export const SETS: Readonly<Record<string, string>> = artifacts;
