import { createFileRoute } from "@tanstack/react-router";
import { Simulator } from "../Pages";

export const Route = createFileRoute("/simulator")({
	component: () => (
		<>
			<title>gcsim - simulator</title>
			<Simulator />
		</>
	),
});
