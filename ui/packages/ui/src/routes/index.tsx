import { createFileRoute } from "@tanstack/react-router";
import { Dash } from "../Pages";

export const Route = createFileRoute("/")({
	component: () => (
		<>
			<title>gcsim - simulation impact</title>
			<Dash />
		</>
	),
});
