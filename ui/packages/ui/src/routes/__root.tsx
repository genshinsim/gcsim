import { ExecutorProvider, useExecutor } from "@gcsim/components";
import type { Executor, ExecutorSupplier } from "@gcsim/executors";
import { Toaster } from "@gcsim/primitives";
import {
	createRootRouteWithContext,
	Outlet,
	useLocation,
	useNavigate,
} from "@tanstack/react-router";
import { useCallback, useEffect, useRef } from "react";
import { useTranslation } from "react-i18next";
import { Footer, Nav } from "../Sectioning";
import { lastRunStore } from "../Stores/lastRun";

export type AppContext = {
	exec: ExecutorSupplier<Executor>;
	mode: string;
	gitCommit: string;
};

export const Route = createRootRouteWithContext<AppContext>()({
	component: Root,
	notFoundComponent: NotFound,
});

const runStore = lastRunStore(localStorage);

function Root() {
	const { exec } = Route.useRouteContext();
	const navigate = useNavigate();
	const navigateOnRun = useCallback(() => navigate({ to: "/web" }), [navigate]);
	return (
		<ExecutorProvider
			exec={exec}
			navigateOnRun={navigateOnRun}
			store={runStore}
		>
			<Layout />
		</ExecutorProvider>
	);
}

function Layout() {
	const content = useRef<HTMLDivElement>(null);
	const href = useLocation({ select: (l) => l.href });
	const pathname = useLocation({ select: (l) => l.pathname });

	// every time you change location, scroll to top of page. This is necessary since the outer
	// content div will never rerender through the entire lifespan of the app and will always retain
	// its scroll position.
	// biome-ignore lint/correctness/useExhaustiveDependencies: href is the trigger, not a value used in the effect body
	useEffect(() => {
		content.current?.scrollTo(0, 0);
	}, [href]);

	const { cancel } = useExecutor();
	const prevPathname = useRef(pathname);
	useEffect(() => {
		if (prevPathname.current === "/web" && pathname !== "/web") {
			cancel();
		}
		prevPathname.current = pathname;
	}, [pathname, cancel]);

	return (
		<div className="h-screen flex flex-col">
			<Toaster position="top-right" theme="dark" />
			<Nav />
			<div
				ref={content}
				className="flex flex-col flex-auto overflow-y-scroll overflow-x-clip"
			>
				<Outlet />
				<Footer />
			</div>
		</div>
	);
}

function NotFound() {
	const { t } = useTranslation();
	return (
		<>
			<title>gcsim - simulation impact</title>
			<div className="m-2 text-center">{t("src.this_page_is")}</div>
		</>
	);
}
