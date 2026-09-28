import type { Executor, ExecutorSupplier } from "@gcsim/executors";
import { createRouter, RouterProvider } from "@tanstack/react-router";
import type { ReactNode } from "react";
import type { AppContext } from "./routes/__root";
import { routeTree } from "./routeTree.gen";
import { AppStateProvider } from "./Stores/AppState";

import "@gcsim/components/src/index.css";
import "./index.css";

// Two things must always be supplied to the UI for it to work
//  1. ExecutorSupplier (Should rarely/never change. Must be react safe)
//  2. ExecutorSettings (passed as children, options and executor state management owned by app)
//
// We use a supplier to give the owning app the opporunity to add their own logic every time we
// try to access the executor. This means they can defer creation to only when it is used
// (improving performance and UX), auto refresh/recreate it if it detects a stale state, or change
// what pool instance is in use
//
// ExecutorSettings is a dialog which will be added to the DOM as part of the main content
// so that it is always loaded in. This is how the app can supply state and decide how it wants to
// construct and configure the executors (and which executors are available to use).
type UIProps = {
	exec: ExecutorSupplier<Executor>;
	children: ReactNode;
	mode: string;
	gitCommit: string;
};

const router = createRouter({
	routeTree,
	context: {} as AppContext,
});

declare module "@tanstack/react-router" {
	interface Register {
		router: typeof router;
	}
}

export const UI = ({ exec, children, mode, gitCommit }: UIProps) => {
	return (
		<AppStateProvider>
			<RouterProvider
				router={router}
				context={{ exec, settings: children, mode, gitCommit }}
			/>
		</AppStateProvider>
	);
};
