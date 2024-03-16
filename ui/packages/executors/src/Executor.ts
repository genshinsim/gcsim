import { ParsedResult, Sample, SimResults } from "@gcsim/types";

export interface Executor {
  ready(): Promise<boolean>;
  running(): boolean;
  validate(cfg: string): Promise<ParsedResult>;
  sample(cfg: string, seed: string): Promise<Sample>;
  run(cfg: string, updateResult: (result: SimResults, hash: string) => void): Promise<boolean | void>;
  cancel(): void;
  buildInfo(): { hash: string; date: string };
}