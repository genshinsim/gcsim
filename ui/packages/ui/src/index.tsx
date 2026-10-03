import type { Executor, ExecutorSupplier } from "@gcsim/executors";
import {
	createRouter,
	parseSearchWith,
	RouterProvider,
	stringifySearchWith,
} from "@tanstack/react-router";
import { type ReactNode, useEffect, useMemo, useRef } from "react";
import type { AppContext } from "./routes/__root";
import { routeTree } from "./routeTree.gen";
import { AppStateProvider } from "./stores/AppState";
import { type ExecutorKind, HostProvider } from "./stores/host";

import "@gcsim/components/src/index.css";
import "./index.css";

// The owning app supplies the executor and its settings:
//  1. ExecutorSupplier (Should rarely/never change. Must be react safe)
//  2. executorSettings + executorKind (options and executor state management owned by app)
//
// We use a supplier to give the owning app the opporunity to add their own logic every time we
// try to access the executor. This means they can defer creation to only when it is used
// (improving performance and UX), auto refresh/recreate it if it detects a stale state, or change
// what pool instance is in use
//
// executorSettings is rendered on the /settings page. This is how the app can supply state and
// decide how it wants to construct and configure the executors (and which executors are available
// to use).
type UIProps = {
	exec: ExecutorSupplier<Executor>;
	executorSettings: ReactNode;
	executorKind: ExecutorKind;
	mode: string;
	gitCommit: string;
};
const router = createRouter({
	routeTree,
	context: {} as AppContext,
	// flat values, so ?seed=123 stays unquoted and a 64-bit seed keeps every digit
	parseSearch: parseSearchWith((value) => value),
	stringifySearch: stringifySearchWith(JSON.stringify),
});

declare module "@tanstack/react-router" {
	interface Register {
		router: typeof router;
	}
}

export const UI = ({
	exec,
	executorSettings,
	executorKind,
	mode,
	gitCommit,
}: UIProps) => {
	const context = useMemo(
		() => ({ exec, mode, gitCommit }),
		[exec, mode, gitCommit],
	);
	const lastContext = useRef(context);
	useEffect(() => {
		if (lastContext.current === context) return;
		lastContext.current = context;
		router.update({ ...router.options, context });
		router.invalidate();
	}, [context]);
	const host = useMemo(
		() => ({ executorSettings, executorKind }),
		[executorSettings, executorKind],
	);

	return (
		<AppStateProvider>
			<HostProvider value={host}>
				<RouterProvider router={router} context={context} />
			</HostProvider>
		</AppStateProvider>
	);
};
