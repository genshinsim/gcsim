import { createFileRoute, redirect } from "@tanstack/react-router";

export const Route = createFileRoute("/s/$id")({
	beforeLoad: ({ params }) => {
		throw redirect({ to: "/sh/$id", params, hash: true, replace: true });
	},
});
