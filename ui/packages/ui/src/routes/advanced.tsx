import { createFileRoute, redirect } from "@tanstack/react-router";

export const Route = createFileRoute("/advanced")({
	beforeLoad: () => {
		throw redirect({ to: "/simulator", replace: true });
	},
});
