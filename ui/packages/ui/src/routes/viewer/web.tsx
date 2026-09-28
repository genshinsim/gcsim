import { createFileRoute, redirect } from "@tanstack/react-router";

export const Route = createFileRoute("/viewer/web")({
	beforeLoad: () => {
		throw redirect({ to: "/web", replace: true });
	},
});
