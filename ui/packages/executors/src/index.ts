import { Executor } from "./Executor";

export type { Executor };
export type ExecutorSupplier = () => Executor;

export { WasmExecutor } from "./WasmExecutor";