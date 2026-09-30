import { createFileRoute } from "@tanstack/react-router";
import { Simulator } from "../features/simulator/Simulator";
import { validateSimulatorSearch } from "../features/simulator/search";

export const Route = createFileRoute("/simulator")({
	validateSearch: validateSimulatorSearch,
	component: () => (
		<>
			<title>gcsim - simulator</title>
			<Simulator />
		</>
	),
});
