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
import { type ReactNode, useEffect, useRef } from "react";
import { Helmet } from "react-helmet";
import { useTranslation } from "react-i18next";
import { Provider } from "react-redux";
import {
	BrowserRouter,
	Navigate,
	Route,
	Routes,
	useLocation,
	useParams,
} from "react-router-dom";
import {
	Dash,
	DBViewer,
	DiscordCallback,
	LocalSample,
	LocalViewer,
	PageUserAccount,
	ShareViewer,
	Simulator,
	UploadSample,
	WebViewer,
} from "./Pages";
import { Footer, Nav } from "./Sectioning";
import { appActions } from "./Stores/appSlice";
import {
	type RootState,
	store,
	useAppDispatch,
	useAppSelector,
} from "./Stores/store";

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

export const UI = (props: UIProps) => {
	return (
		<BrowserRouter>
			<Provider store={store}>
				<Main {...props} />
			</Provider>
		</BrowserRouter>
	);
};

function RedirectDB() {
	window.location.replace("https://db.gcsim.app");
	return (
		<div>
			Please visit the new db at{" "}
			<a href="https://db.gcsim.app">https://db.gcsim.app</a>
		</div>
	);
}

// TODO: Move to its own file?
// TODO: Add tabs for better settings management + extensibility?
const ExecutorSettings = ({ children }: { children: ReactNode }) => {
	const { t } = useTranslation();
	const dispatch = useAppDispatch();
	const { isOpen, sampleOnLoad } = useAppSelector((state: RootState) => {
		return {
			isOpen: state.app.isSettingsOpen,
			sampleOnLoad: state.app.sampleOnLoad,
		};
	});

	return (
		<Dialog
			open={isOpen}
			onOpenChange={(open) =>
				!open && dispatch(appActions.setSettingsOpen(false))
			}
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
						onCheckedChange={() =>
							dispatch(appActions.setSampleOnLoad(!sampleOnLoad))
						}
					/>
					<Label htmlFor="sample-on-load">{t("simple.generate_sample")}</Label>
				</div>
			</DialogContent>
		</Dialog>
	);
};

const viewerPaths = ["/web", "/local", "/sh/", "/db/"];

type ShareRouteProps = {
	exec: ExecutorSupplier<Executor>;
	gitCommit: string;
	mode: string;
};

function ShareViewerRoute({ exec, gitCommit, mode }: ShareRouteProps) {
	const { id } = useParams();
	useEffect(() => {
		document.title = "gcsim sh - " + id;
	}, [id]);
	return <ShareViewer exec={exec} id={id} gitCommit={gitCommit} mode={mode} />;
}

function DBViewerRoute({ exec, gitCommit, mode }: ShareRouteProps) {
	const { id } = useParams();
	useEffect(() => {
		document.title = "gcsim db - " + id;
	}, [id]);
	return <DBViewer exec={exec} id={id} gitCommit={gitCommit} mode={mode} />;
}

function RedirectToShare() {
	const { id } = useParams();
	return <Navigate to={"/sh/" + id} replace />;
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
function movedOffViewer(location: any, prevLocation: any): boolean {
	let prevWasViewer = false;
	let destIsViewer = false;
	for (let i = 0; i < viewerPaths.length; i++) {
		prevWasViewer =
			prevWasViewer || prevLocation.current.pathname.startsWith(viewerPaths[i]);
		destIsViewer = destIsViewer || location.pathname.startsWith(viewerPaths[i]);
	}
	return prevWasViewer && !destIsViewer;
}

const Main = ({ exec, children, gitCommit, mode }: UIProps) => {
	const { t } = useTranslation();
	const content = useRef<HTMLDivElement>(null);
	const location = useLocation();

	// every time you change location, scroll to top of page. This is necessary since the outer
	// content div will never rerender through the entire lifespan of the app and will always retain
	// its scroll position.
	// biome-ignore lint/correctness/useExhaustiveDependencies: location is the trigger, not a value used in the effect body
	useEffect(() => {
		content.current?.scrollTo(0, 0);
	}, [location]);

	// cancel the run every time we navigate away from the web viewer page
	const prevLocation = useRef(location);
	useEffect(() => {
		if (
			prevLocation.current !== location &&
			movedOffViewer(location, prevLocation) &&
			exec().running()
		) {
			exec().cancel();
		}
		prevLocation.current = location;
	}, [location, exec]);

	return (
		<div className="h-screen flex flex-col">
			<Toaster position="top-right" theme="dark" />
			<Nav />
			<div
				ref={content}
				className="flex flex-col flex-auto overflow-y-scroll overflow-x-clip"
			>
				<Routes>
					<Route
						path="/"
						element={
							<>
								<Helmet>
									<title>gcsim - simulation impact</title>
								</Helmet>
								<Dash />
							</>
						}
					/>

					{/* Simulator */}
					<Route
						path="/simulator"
						element={
							<>
								<Helmet>
									<title>gcsim - simulator</title>
								</Helmet>
								<Simulator exec={exec} />
							</>
						}
					/>

					{/* Viewer Routes */}
					<Route
						path="/web/*"
						element={
							<>
								<Helmet>
									<title>gcsim - viewer</title>
								</Helmet>
								<WebViewer exec={exec} gitCommit={gitCommit} mode={mode} />
							</>
						}
					/>
					<Route
						path="/local/*"
						element={
							<>
								<Helmet>
									<title>gcsim - local viewer</title>
								</Helmet>
								<LocalViewer exec={exec} gitCommit={gitCommit} mode={mode} />
							</>
						}
					/>
					<Route
						path="/sh/:id"
						element={
							<ShareViewerRoute exec={exec} gitCommit={gitCommit} mode={mode} />
						}
					/>
					<Route
						path="/db/:id"
						element={
							<DBViewerRoute exec={exec} gitCommit={gitCommit} mode={mode} />
						}
					/>

					{/* Sample Routes */}
					<Route
						path="/sample/upload"
						element={
							<>
								<Helmet>
									<title>gcsim - sample</title>
								</Helmet>
								<UploadSample />
							</>
						}
					/>
					<Route
						path="/sample/local"
						element={
							<>
								<Helmet>
									<title>gcsim - local sample</title>
								</Helmet>
								<LocalSample />
							</>
						}
					/>

					{/* Redirects */}
					<Route path="/v3/viewer/share/:id" element={<RedirectToShare />} />
					<Route path="/viewer/share/:id" element={<RedirectToShare />} />
					<Route path="/s/:id" element={<RedirectToShare />} />
					<Route path="/viewer/web" element={<Navigate to="/web" replace />} />
					<Route
						path="/viewer/local"
						element={<Navigate to="/local" replace />}
					/>
					<Route
						path="/simple"
						element={<Navigate to="/simulator" replace />}
					/>
					<Route
						path="/advanced"
						element={<Navigate to="/simulator" replace />}
					/>
					<Route
						path="/viewer"
						element={<Navigate to="/simulator" replace />}
					/>

					{/* DB & Account */}
					<Route path="/db" element={<RedirectDB />} />
					<Route
						path="/account"
						element={
							<>
								<Helmet>
									<title>gcsim - account</title>
								</Helmet>
								<PageUserAccount />
							</>
						}
					/>
					<Route path="/auth/discord" element={<DiscordCallback />} />

					{/* Default (404 case) */}
					<Route
						path="*"
						element={
							<>
								<Helmet>
									<title>gcsim - simulation impact</title>
								</Helmet>
								<div className="m-2 text-center">{t("src.this_page_is")}</div>
							</>
						}
					/>
				</Routes>
				<Footer />
				<ExecutorSettings>{children}</ExecutorSettings>
			</div>
		</div>
	);
};
