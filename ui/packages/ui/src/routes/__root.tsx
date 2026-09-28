import { ExecutorProvider, useExecutor } from "@gcsim/components";
import type { Executor, ExecutorSupplier } from "@gcsim/executors";
import {
	Dialog,
	DialogContent,
	DialogHeader,
	DialogTitle,
	Label,
	Switch as SwitchInput,
	Toaster,
} from "@gcsim/primitives";
import {
	createRootRouteWithContext,
	Outlet,
	useLocation,
	useNavigate,
} from "@tanstack/react-router";
import { type ReactNode, useCallback, useEffect, useRef } from "react";
import { useTranslation } from "react-i18next";
import { Footer, Nav } from "../Sectioning";
import { usePrefs } from "../Stores/AppState";
import { lastRunStore } from "../Stores/lastRun";

export type AppContext = {
	exec: ExecutorSupplier<Executor>;
	settings: ReactNode;
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
	const { settings } = Route.useRouteContext();
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
				<ExecutorSettings>{settings}</ExecutorSettings>
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
