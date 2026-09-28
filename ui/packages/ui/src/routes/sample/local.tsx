import { createFileRoute } from "@tanstack/react-router";
import { LocalSample } from "../../Pages";

export const Route = createFileRoute("/sample/local")({
	component: () => (
		<>
			<title>gcsim - local sample</title>
			<LocalSample />
		</>
	),
});
