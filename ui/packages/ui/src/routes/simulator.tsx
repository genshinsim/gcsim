import { createFileRoute } from "@tanstack/react-router";
import { Simulator } from "../features/simulator/Simulator";

export const Route = createFileRoute("/simulator")({
	component: () => (
		<>
			<title>gcsim - simulator</title>
			<Simulator />
		</>
	),
});
