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

const ExecutorContext = React.createContext<ExecutorContextValue | null>(null);
const RunResultContext = React.createContext<RunResult | null>(null);

export interface ExecutorProviderProps {
	exec: ExecutorSupplier<Executor>;
	onResult?: (result: model.SimulationResult, hash: string) => void;
	navigateOnRun?: () => void;
	children: React.ReactNode;
}

const noop = () => {};

export function ExecutorProvider({
	exec,
	onResult,
	navigateOnRun,
	children,
}: ExecutorProviderProps) {
	const [isReady, setReady] = React.useState(false);
	const [busy, setBusy] = React.useState(false);
	const [runResult, setRunResult] = React.useState<RunResult>(emptyRunResult);

	const onResultRef = React.useRef(onResult);
	onResultRef.current = onResult;
	const navigateOnRunRef = React.useRef(navigateOnRun);
	navigateOnRunRef.current = navigateOnRun;
	const cancelSinkRef = React.useRef(noop);

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
			cancelSinkRef.current();
		};
	}, [exec]);

	const run = React.useCallback(
		(config: string) => {
			const executor = exec();
			if (executor.running()) {
				return;
			}
			cancelSinkRef.current();
			setRunResult(emptyRunResult);
			executor.validate(config).then((validated) => {
				if (
					(validated.errors && validated.errors.length > 0) ||
					executor.running()
				) {
					return;
				}
				const sink = throttle(
					(result: model.SimulationResult, hash: string) => {
						onResultRef.current?.(result, hash);
						React.startTransition(() =>
							setRunResult((prev) => ({ ...prev, result, hash })),
						);
					},
					RESULT_THROTTLE_MS,
					{ leading: true, trailing: true },
				);
				cancelSinkRef.current = () => sink.cancel();
				setRunResult({ ...emptyRunResult, config });
				executor.run(config, sink).catch((err: unknown) => {
					setRunResult((prev) => ({ ...prev, error: asError(err) }));
				});
				navigateOnRunRef.current?.();
			}, noop);
		},
		[exec],
	);

	const value = React.useMemo<ExecutorContextValue>(
		() => ({ exec, isReady, busy, run }),
		[exec, isReady, busy, run],
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
