import { createFileRoute } from "@tanstack/react-router";
import { Dash } from "../features/dashboard/Dash";

export const Route = createFileRoute("/")({
	component: () => (
		<>
			<title>gcsim - simulation impact</title>
			<Dash />
		</>
	),
});
