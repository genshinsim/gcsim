import type { Executor, ExecutorSupplier, model } from "@gcsim/types";
import { throttle } from "lodash-es";
import React from "react";
import { asError } from "./asError";

const READY_POLL_MS = 250;
const BUSY_POLL_MS = 50;
const RESULT_THROTTLE_MS = 100;

export interface ExecutorContextValue {
	exec: ExecutorSupplier<Executor>;
	isReady: boolean;
	busy: boolean;
	run: (config: string) => void;
	cancel: () => void;
}

export interface RunResult {
	result: model.SimulationResult | null;
	hash: string | null;
	config: string | null;
	error: string | null;
}

const emptyRunResult: RunResult = {
	result: null,
	hash: null,
	config: null,
	error: null,
};

export interface SavedRun {
	result: model.SimulationResult;
	hash: string | null;
}

export interface SavedRunStore {
	load: () => SavedRun | null;
	save: (run: SavedRun) => void;
}

function initialRunResult(store?: SavedRunStore): RunResult {
	const saved = store?.load();
	return saved == null ? emptyRunResult : { ...emptyRunResult, ...saved };
}

const ExecutorContext = React.createContext<ExecutorContextValue | null>(null);
const RunResultContext = React.createContext<RunResult | null>(null);

export interface ExecutorProviderProps {
	exec: ExecutorSupplier<Executor>;
	navigateOnRun?: () => void;
	store?: SavedRunStore;
	children: React.ReactNode;
}

const noop = () => {};

export function ExecutorProvider({
	exec,
	navigateOnRun,
	store,
	children,
}: ExecutorProviderProps) {
	const [isReady, setReady] = React.useState(false);
	const [busy, setBusy] = React.useState(false);
	const [runResult, setRunResult] = React.useState(() =>
		initialRunResult(store),
	);

	const storeRef = React.useRef(store);
	storeRef.current = store;
	const navigateOnRunRef = React.useRef(navigateOnRun);
	navigateOnRunRef.current = navigateOnRun;
	const disarmRef = React.useRef(noop);

	React.useEffect(() => {
		let active = true;
		const pollReady = () =>
			exec()
				.ready()
				.then((res) => {
					if (active) {
						setReady(res);
					}
				});
		const pollBusy = () => {
			if (active) {
				setBusy(exec().running());
			}
		};
		pollReady();
		pollBusy();
		const readyInterval = setInterval(pollReady, READY_POLL_MS);
		const busyInterval = setInterval(pollBusy, BUSY_POLL_MS);
		return () => {
			active = false;
			clearInterval(readyInterval);
			clearInterval(busyInterval);
			disarmRef.current();
		};
	}, [exec]);

	const run = React.useCallback(
		(config: string) => {
			const executor = exec();
			if (executor.running()) {
				return;
			}
			disarmRef.current();
			setRunResult(emptyRunResult);
			executor.validate(config).then((validated) => {
				if (
					(validated.errors && validated.errors.length > 0) ||
					executor.running()
				) {
					return;
				}
				let armed = true;
				let last: SavedRun | null = null;
				const apply = throttle(
					(result: model.SimulationResult, hash: string) => {
						React.startTransition(() =>
							setRunResult((prev) => ({ ...prev, result, hash })),
						);
					},
					RESULT_THROTTLE_MS,
					{ leading: true, trailing: true },
				);
				const sink = (result: model.SimulationResult, hash: string) => {
					if (armed) {
						last = { result, hash };
						apply(result, hash);
					}
				};
				disarmRef.current = () => {
					armed = false;
					apply.cancel();
				};
				setRunResult({ ...emptyRunResult, config });
				executor.run(config, sink).then(
					() => {
						if (armed && last != null) {
							storeRef.current?.save(last);
						}
					},
					(err: unknown) => {
						if (armed) {
							setRunResult((prev) => ({ ...prev, error: asError(err) }));
						}
					},
				);
				navigateOnRunRef.current?.();
			}, noop);
		},
		[exec],
	);

	const cancel = React.useCallback(() => {
		disarmRef.current();
		const executor = exec();
		if (executor.running()) {
			executor.cancel();
		}
	}, [exec]);

	const value = React.useMemo<ExecutorContextValue>(
		() => ({ exec, isReady, busy, run, cancel }),
		[exec, isReady, busy, run, cancel],
	);

	return (
		<ExecutorContext.Provider value={value}>
			<RunResultContext.Provider value={runResult}>
				{children}
			</RunResultContext.Provider>
		</ExecutorContext.Provider>
	);
}

export function useExecutor(): ExecutorContextValue {
	const ctx = React.useContext(ExecutorContext);
	if (ctx == null) {
		throw new Error("useExecutor must be used within an ExecutorProvider");
	}
	return ctx;
}

export function useRunResult(): RunResult {
	const ctx = React.useContext(RunResultContext);
	if (ctx == null) {
		throw new Error("useRunResult must be used within an ExecutorProvider");
	}
	return ctx;
}
