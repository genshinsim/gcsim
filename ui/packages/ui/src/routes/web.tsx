import { createFileRoute, redirect } from "@tanstack/react-router";
import { WebViewer } from "../Pages";
import { legacyHashSearch, validateViewerSearch } from "../Pages/Viewer/search";

export const Route = createFileRoute("/web")({
	validateSearch: validateViewerSearch,
	beforeLoad: ({ location }) => {
		const search = legacyHashSearch(location.hash);
		if (search != null) {
			throw redirect({ to: "/web", search, replace: true });
		}
	},
	component: Page,
});

function Page() {
	const { gitCommit, mode } = Route.useRouteContext();
	return (
		<>
			<title>gcsim - viewer</title>
			<WebViewer gitCommit={gitCommit} mode={mode} />
		</>
	);
}
