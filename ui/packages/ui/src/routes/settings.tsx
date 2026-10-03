import { createFileRoute } from "@tanstack/react-router";
import { Settings } from "../features/settings/Settings";

export const Route = createFileRoute("/settings")({
	component: () => (
		<>
			<title>gcsim - settings</title>
			<Settings />
		</>
	),
});
