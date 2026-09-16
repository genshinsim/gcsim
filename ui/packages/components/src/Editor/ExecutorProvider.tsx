import type { Executor, ExecutorSupplier } from "@gcsim/types";
import React from "react";

const READY_POLL_MS = 250;
const RUNNING_POLL_MS = 50;

export interface ExecutorContextValue {
	exec: ExecutorSupplier<Executor>;
	isReady: boolean;
	running: boolean;
}

const ExecutorContext = React.createContext<ExecutorContextValue | null>(null);

export interface ExecutorProviderProps {
	exec: ExecutorSupplier<Executor>;
	children: React.ReactNode;
}

export function ExecutorProvider({ exec, children }: ExecutorProviderProps) {
	const [isReady, setReady] = React.useState(false);
	const [running, setRunning] = React.useState(false);

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

	const value = React.useMemo<ExecutorContextValue>(
		() => ({ exec, isReady, running }),
		[exec, isReady, running],
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
