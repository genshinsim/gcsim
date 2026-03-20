declare global {
  function importScripts(script: string): void;

  class Go {
    argv: string[];
    env: { [envKey: string]: string };
    exit: (code: number) => void;
    importObject: WebAssembly.Imports;
    exited: boolean;
    mem: DataView;
    run(instance: WebAssembly.Instance): Promise<void>;
  }

  // Helper functions
  function sample(cfg: string, seed: string): string;
  function validateConfig(cfg: string): string;

  // Aggregator functions
  function initializeAggregator(cfg: string): string;
  function aggregate(result: Uint8Array): string | null;
  function flush(): string;

  // Worker functions
  function initializeWorker(cfg: string): string | null;
  function simulate(): Uint8Array | string;
}

export {};
