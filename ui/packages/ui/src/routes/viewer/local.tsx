import { createFileRoute, redirect } from "@tanstack/react-router";

export const Route = createFileRoute("/viewer/local")({
	beforeLoad: () => {
		throw redirect({ to: "/local", hash: true, replace: true });
	},
});
