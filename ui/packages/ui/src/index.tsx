import type { Executor, ExecutorSupplier } from "@gcsim/executors";
import {
	Dialog,
	DialogContent,
	DialogHeader,
	DialogTitle,
	Label,
	Switch as SwitchInput,
} from "@gcsim/primitives";
import {
	createRouter,
	parseSearchWith,
	RouterProvider,
	stringifySearchWith,
} from "@tanstack/react-router";
import { type ReactNode, useEffect, useMemo } from "react";
import { useTranslation } from "react-i18next";
import type { AppContext } from "./routes/__root";
import { routeTree } from "./routeTree.gen";
import { AppStateProvider, usePrefs } from "./stores/AppState";

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
	// flat values, so ?seed=123 stays unquoted and a 64-bit seed keeps every digit
	parseSearch: parseSearchWith((value) => value),
	stringifySearch: stringifySearchWith(JSON.stringify),
});

declare module "@tanstack/react-router" {
	interface Register {
		router: typeof router;
	}
}

export const UI = ({ exec, children, mode, gitCommit }: UIProps) => {
	const context = useMemo(
		() => ({ exec, mode, gitCommit }),
		[exec, mode, gitCommit],
	);
	useEffect(() => {
		router.update({ ...router.options, context });
		router.invalidate();
	}, [context]);

	return (
		<AppStateProvider>
			<RouterProvider router={router} context={context} />
			<ExecutorSettings>{children}</ExecutorSettings>
		</AppStateProvider>
	);
};

// TODO: Add tabs for better settings management + extensibility?
const ExecutorSettings = ({ children }: { children: ReactNode }) => {
	const { t } = useTranslation();
	const { settingsOpen, setSettingsOpen, sampleOnLoad, setSampleOnLoad } =
		usePrefs();

	return (
		<Dialog
			open={settingsOpen}
			onOpenChange={(open) => !open && setSettingsOpen(false)}
		>
			<DialogContent>
				<DialogHeader>
					<DialogTitle>{t("simple.settings")}</DialogTitle>
				</DialogHeader>
				{children}
				<div className="flex items-center gap-2 pt-5">
					<SwitchInput
						id="sample-on-load"
						checked={sampleOnLoad}
						onCheckedChange={() => setSampleOnLoad(!sampleOnLoad)}
					/>
					<Label htmlFor="sample-on-load">{t("simple.generate_sample")}</Label>
				</div>
			</DialogContent>
		</Dialog>
	);
};
