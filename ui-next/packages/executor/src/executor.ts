import type { Sim } from "@gcsim/types";

export interface Executor {
  ready(): Promise<boolean>;
  running(): boolean;
  validate(cfg: string): Promise<Sim.ParsedResult>;
  sample(cfg: string, seed: string): Promise<Sim.Sample>;
  run(
    cfg: string,
    updateResult: (result: Sim.SimResults, hash: string) => void,
  ): Promise<boolean | void>;
  cancel(): void;
  buildInfo(): { hash: string; date: string };
}

export type ExecutorSupplier<T extends Executor> = () => T;
