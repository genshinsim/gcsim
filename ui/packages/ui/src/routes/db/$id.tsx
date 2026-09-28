import { createFileRoute } from "@tanstack/react-router";
import { DBViewer } from "../../Pages";

export const Route = createFileRoute("/db/$id")({
	component: Page,
});

function Page() {
	const { id } = Route.useParams();
	const { gitCommit, mode } = Route.useRouteContext();
	return (
		<>
			<title>{`gcsim db - ${id}`}</title>
			<DBViewer id={id} gitCommit={gitCommit} mode={mode} />
		</>
	);
}
