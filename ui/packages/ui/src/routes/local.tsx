import { createFileRoute } from "@tanstack/react-router";
import { LocalViewer } from "../Pages";

export const Route = createFileRoute("/local")({
	component: Page,
});

function Page() {
	const { gitCommit, mode } = Route.useRouteContext();
	return (
		<>
			<title>gcsim - local viewer</title>
			<LocalViewer gitCommit={gitCommit} mode={mode} />
		</>
	);
}
