import type { model } from "./generated";
import type { ParsedResult, Sample } from "./sim";

export interface Executor {
	ready(): Promise<boolean>;
	running(): boolean;
	validate(cfg: string): Promise<ParsedResult>;
	sample(cfg: string, seed: string): Promise<Sample>;
	run(
		cfg: string,
		updateResult: (
			result: model.SimulationResult,
			raw: string,
			hash: string,
		) => void,
	): Promise<boolean | void>;
	cancel(): void;
	buildInfo(): { hash: string; date: string };
}

export type ExecutorSupplier<T extends Executor> = () => T;
