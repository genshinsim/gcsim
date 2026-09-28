import { createFileRoute, redirect } from "@tanstack/react-router";
import { LoadedViewer, type LoadedViewerProps, loadResult } from "../Pages";
import { legacyHashSearch, validateViewerSearch } from "../Pages/Viewer/search";

export const Route = createFileRoute("/local")({
	validateSearch: validateViewerSearch,
	beforeLoad: ({ location }) => {
		const search = legacyHashSearch(location.hash);
		if (search != null) {
			throw redirect({ to: "/local", search, replace: true });
		}
	},
	loader: () => loadResult("http://127.0.0.1:8381/data"),
	component: () => <Page result={Route.useLoaderData()} />,
	pendingComponent: () => <Page />,
	errorComponent: ({ error }) => <Page error={(error as Error).message} />,
});

function Page(props: Pick<LoadedViewerProps, "result" | "error">) {
	const { gitCommit, mode } = Route.useRouteContext();
	return (
		<>
			<title>gcsim - local viewer</title>
			<LoadedViewer {...props} gitCommit={gitCommit} mode={mode} />
		</>
	);
}
