export interface Sample {
	config?: string;
	initial_character?: string;
	character_details?: Character[];
	target_details?: unknown[];
	seed?: string;
	logs?: LogDetails[];
}

export interface Character {
	name: string;
	level: number;
	element: string;
	max_level: number;
	cons: number;
	weapon: Weapon;
	talents: Talent;
	stats: number[];
	snapshot: number[];
	sets: Set;
	date_added?: string;
	enka_build_name?: string;
	source?: string; // allow for both GOOD and enka import
}

export interface Talent {
	attack: number;
	skill: number;
	burst: number;
}

export interface Set {
	[key: string]: number;
}

export interface Weapon {
	name: string;
	refine: number;
	level: number;
	max_level: number;
}

export type LogDetails = {
	char_index: number;
	ended: number;
	event: string;
	frame: number;
	msg: string;
	// eslint-disable-next-line @typescript-eslint/no-explicit-any
	logs: { [key in string]: any };
	ordering?: { [key: string]: number };
};

export type StatusType = "idle" | "loading" | "done" | "error";

export interface ParsedResult {
	characters: ParsedCharacterProfile[];
	errors: string[];
	player_initial_pos: { x: number; y: number; r: number };
}

export interface ParsedCharacterProfile {
	base: Base;
	weapon: Weapon;
	talents: Talent;
	stats: number[];
	sets: Set;
}

export interface Base {
	key: string;
	name: string;
	element: string;
	level: number;
	max_level: number;
	base_hp: number;
	base_atk: number;
	base_def: number;
	cons: number;
	start_hp: number;
}

export interface ParsedWeapon {
	name: string;
	key: string;
	refine: number;
	level: number;
	max_level: number;
	base_atk: number;
}
