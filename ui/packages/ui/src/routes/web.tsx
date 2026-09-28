import { createFileRoute } from "@tanstack/react-router";
import { WebViewer } from "../Pages";

export const Route = createFileRoute("/web")({
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
