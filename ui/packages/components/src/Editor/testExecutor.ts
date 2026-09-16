import type {
	Executor,
	ParsedCharacterProfile,
	ParsedResult,
	Sample,
	SimResults,
} from "@gcsim/types";
import { vi } from "vitest";

export function makeProfile(key: string): ParsedCharacterProfile {
	return {
		base: {
			key,
			name: key,
			element: "pyro",
			level: 90,
			max_level: 90,
			base_hp: 1,
			base_atk: 1,
			base_def: 1,
			cons: 0,
			start_hp: 1,
		},
		weapon: { name: "w", refine: 1, level: 90, max_level: 90 },
		talents: { attack: 9, skill: 9, burst: 9 },
		stats: [1, 2, 3],
		sets: {},
	};
}

export function parsedResult(
	keys: string[],
	errors: string[] = [],
): ParsedResult {
	return {
		characters: keys.map(makeProfile),
		errors,
		player_initial_pos: { x: 0, y: 0, r: 0 },
	};
}

export interface FakeExecutorOptions {
	ready?: boolean;
	running?: boolean;
	validate?: (cfg: string) => Promise<ParsedResult>;
}

export function makeExecutor(opts: FakeExecutorOptions = {}) {
	let isRunning = opts.running ?? false;
	const validate = vi.fn(
		opts.validate ??
			((_cfg: string) => Promise.resolve(parsedResult(["amber"]))),
	);
	const run = vi.fn(
		(_cfg: string, _sink: (r: SimResults, hash: string) => void) => {
			isRunning = true;
			return Promise.resolve(true);
		},
	);
	const executor: Executor = {
		ready: () => Promise.resolve(opts.ready ?? true),
		running: () => isRunning,
		validate,
		sample: (_cfg: string, _seed: string) => Promise.resolve({} as Sample),
		run,
		cancel: () => {
			isRunning = false;
		},
		buildInfo: () => ({ hash: "", date: "" }),
	};
	return {
		executor,
		supplier: () => executor,
		validate,
		run,
		setRunning: (v: boolean) => {
			isRunning = v;
		},
	};
}
