import { createFileRoute, redirect } from "@tanstack/react-router";

export const Route = createFileRoute("/viewer/share/$id")({
	beforeLoad: ({ params }) => {
		throw redirect({ to: "/sh/$id", params, replace: true });
	},
});
