import { createFileRoute } from "@tanstack/react-router";
import { UploadSample } from "../../features/sample/UploadSample";

export const Route = createFileRoute("/sample/upload")({
	component: () => (
		<>
			<title>gcsim - sample</title>
			<UploadSample />
		</>
	),
});
