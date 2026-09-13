import type { Executor } from "./Executor";

export type { Executor };
export type ExecutorSupplier<T extends Executor> = () => T;

export { ServerExecutor } from "./ServerExecutor";
export { WasmExecutor } from "./WasmExecutor";
