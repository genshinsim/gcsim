import { LogDetails, ParsedResult, SimResults } from "@gcsim/types";

export interface Executor {
  count(): number;
  setWorkerCount(count: number): void;
  ready(): boolean;
  running(): boolean;
  validate(cfg: string): Promise<ParsedResult>;
  debug(cfg: string, seed: string): Promise<LogDetails[]>;
  run(cfg: string, updateResult: (result: SimResults) => void): Promise<boolean | void>;
  cancel(): void;
  buildInfo(): { hash: string; date: string };
}