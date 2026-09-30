import { ExecutorProvider, useExecutor } from "@gcsim/components";
import type { Executor, ExecutorSupplier } from "@gcsim/executors";
import { Toaster } from "@gcsim/primitives";
import {
	createRootRouteWithContext,
	Outlet,
	useLocation,
	useNavigate,
	useRouter,
} from "@tanstack/react-router";
import { useCallback, useEffect, useRef } from "react";
import { useTranslation } from "react-i18next";
import { Footer, Nav } from "../features/layout";
import { lastRunStore } from "../stores/lastRun";

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
	const router = useRouter();
	const pathname = useLocation({ select: (l) => l.pathname });

	// every time you change location, scroll to top of page. This is necessary since the outer
	// content div will never rerender through the entire lifespan of the app and will always retain
	// its scroll position.
	useEffect(
		() => router.subscribe("onResolved", () => content.current?.scrollTo(0, 0)),
		[router],
	);

	const { cancel } = useExecutor();
	const prevPathname = useRef(pathname);
	useEffect(() => {
		if (prevPathname.current === "/web" && pathname !== "/web") {
			cancel();
		}
		prevPathname.current = pathname;
	}, [pathname, cancel]);

	return (
		<div className="h-dvh flex flex-col pt-[env(safe-area-inset-top)] pl-[env(safe-area-inset-left)] pr-[env(safe-area-inset-right)]">
			<Toaster
				position="top-right"
				theme="dark"
				offset={{
					top: "calc(24px + env(safe-area-inset-top))",
					right: "calc(24px + env(safe-area-inset-right))",
				}}
				mobileOffset={{
					top: "calc(16px + env(safe-area-inset-top))",
					right: "calc(16px + env(safe-area-inset-right))",
					left: "calc(16px + env(safe-area-inset-left))",
				}}
			/>
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
