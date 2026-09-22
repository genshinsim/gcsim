import type {
	Executor,
	ExecutorSupplier,
	SimResults,
} from "@gcsim/types";
import React from "react";

const READY_POLL_MS = 250;
const RUNNING_POLL_MS = 50;

export interface ExecutorContextValue {
	exec: ExecutorSupplier<Executor>;
	isReady: boolean;
	running: boolean;
	run: (config: string) => void;
}

const ExecutorContext = React.createContext<ExecutorContextValue | null>(null);

export interface ExecutorProviderProps {
	exec: ExecutorSupplier<Executor>;
	onResult?: (result: SimResults, hash: string) => void;
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
	const [running, setRunning] = React.useState(false);

	const onResultRef = React.useRef(onResult);
	onResultRef.current = onResult;
	const navigateOnRunRef = React.useRef(navigateOnRun);
	navigateOnRunRef.current = navigateOnRun;

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
		const pollRunning = () => {
			if (active) {
				setRunning(exec().running());
			}
		};
		pollReady();
		pollRunning();
		const readyInterval = setInterval(pollReady, READY_POLL_MS);
		const runningInterval = setInterval(pollRunning, RUNNING_POLL_MS);
		return () => {
			active = false;
			clearInterval(readyInterval);
			clearInterval(runningInterval);
		};
	}, [exec]);

	const run = React.useCallback(
		(config: string) => {
			const executor = exec();
			executor.validate(config).then((result) => {
				if ((result.errors && result.errors.length > 0) || executor.running()) {
					return;
				}
				executor.run(config, onResultRef.current ?? noop);
				navigateOnRunRef.current?.();
			}, noop);
		},
		[exec],
	);

	const value = React.useMemo<ExecutorContextValue>(
		() => ({ exec, isReady, running, run }),
		[exec, isReady, running, run],
	);

	return (
		<ExecutorContext.Provider value={value}>
			{children}
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
