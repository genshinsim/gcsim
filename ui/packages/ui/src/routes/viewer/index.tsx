import { createFileRoute, redirect } from "@tanstack/react-router";

export const Route = createFileRoute("/viewer/")({
	beforeLoad: () => {
		throw redirect({ to: "/simulator", replace: true });
	},
});
