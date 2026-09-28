import { createFileRoute } from "@tanstack/react-router";
import { ShareViewer } from "../../Pages";

export const Route = createFileRoute("/sh/$id")({
	component: Page,
});

function Page() {
	const { id } = Route.useParams();
	const { gitCommit, mode } = Route.useRouteContext();
	return (
		<>
			<title>{`gcsim sh - ${id}`}</title>
			<ShareViewer id={id} gitCommit={gitCommit} mode={mode} />
		</>
	);
}
