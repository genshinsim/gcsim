import { createFileRoute } from "@tanstack/react-router";
import { PageUserAccount } from "../features/user/PageUserAccount";

export const Route = createFileRoute("/account")({
	component: () => (
		<>
			<title>gcsim - account</title>
			<PageUserAccount />
		</>
	),
});
