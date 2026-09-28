import { createFileRoute } from "@tanstack/react-router";
import { UploadSample } from "../../Pages";

export const Route = createFileRoute("/sample/upload")({
	component: () => (
		<>
			<title>gcsim - sample</title>
			<UploadSample />
		</>
	),
});
