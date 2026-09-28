import { createFileRoute, redirect } from "@tanstack/react-router";

export const Route = createFileRoute("/simple")({
	beforeLoad: () => {
		throw redirect({ to: "/simulator", replace: true });
	},
});
