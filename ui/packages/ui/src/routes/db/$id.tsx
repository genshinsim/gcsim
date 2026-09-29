import { createFileRoute, redirect } from "@tanstack/react-router";
import {
	LoadedViewer,
	type LoadedViewerProps,
	loadResult,
} from "../../features/viewer";
import {
	legacyHashSearch,
	validateViewerSearch,
} from "../../features/viewer/search";

export const Route = createFileRoute("/db/$id")({
	validateSearch: validateViewerSearch,
	beforeLoad: ({ location, params }) => {
		const search = legacyHashSearch(location.hash);
		if (search != null) {
			throw redirect({ to: "/db/$id", params, search, replace: true });
		}
	},
	shouldReload: ({ cause }) => cause === "enter",
	loader: ({ params }) => loadResult(`/api/share/db/${params.id}`),
	component: () => <Page result={Route.useLoaderData()} />,
	pendingComponent: () => <Page />,
	errorComponent: ({ error }) => <Page error={(error as Error).message} />,
});

function Page(props: Pick<LoadedViewerProps, "result" | "error">) {
	const { id } = Route.useParams();
	const { gitCommit, mode } = Route.useRouteContext();
	return (
		<>
			<title>{`gcsim db - ${id}`}</title>
			<LoadedViewer {...props} gitCommit={gitCommit} mode={mode} />
		</>
	);
}
