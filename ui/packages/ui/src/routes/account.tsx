import { createFileRoute } from "@tanstack/react-router";
import { PageUserAccount } from "../Pages";

export const Route = createFileRoute("/account")({
	component: () => (
		<>
			<title>gcsim - account</title>
			<PageUserAccount />
		</>
	),
});
