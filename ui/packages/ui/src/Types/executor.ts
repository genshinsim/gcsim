import { ParsedResult, LogDetails, SimResults } from "./sim";

export interface Executor {
  count(): number;
  setWorkerCount(count: number, onReady: (count: number) => void): void;
  ready(): boolean;
  running(): boolean;
  validate(cfg: string): Promise<ParsedResult>;
  debug(cfg: string, seed: string): Promise<LogDetails[]>;
  run(
    cfg: string,
    updateResult: (result: SimResults) => void
  ): Promise<boolean | void>;
  cancel(): void;
  buildInfo(): { hash: string; date: string };
}
