declare global {
  declare function importScripts(script: string);

  declare class Go {
    argv: string[];
    env: { [envKey: string]: string };
    exit: (code: number) => void;
    importObject: WebAssembly.Imports;
    exited: boolean;
    mem: DataView;
    run(instance: WebAssembly.Instance): Promise<void>;
  }
  
  // Helper functions
  declare function buildInfo(): string;
  declare function validateConfig(cfg: string): string;

  // Aggregator functions
  declare function initializeAggregator(cfg: string): string | null;
  declare function aggregate(result: Uint8Array, itr: int): string | null;
  declare function flush(): string;

  // Worker functions
  declare function initializeWorker(cfg: string): string | null;
  declare function simulate(): Uint8Array | string;
}

export {};
