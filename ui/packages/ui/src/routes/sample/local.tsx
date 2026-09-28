import type { Sample } from "@gcsim/types";
import { createFileRoute, useRouter } from "@tanstack/react-router";
import axios from "axios";
import { SamplePage } from "../../Pages";

export const Route = createFileRoute("/sample/local")({
	shouldReload: ({ cause }) => cause === "enter",
	loader: async () =>
		(
			await axios.get<Sample>("http://127.0.0.1:8381/sample", {
				timeout: 30000,
			})
		).data,
	component: () => <Page sample={Route.useLoaderData()} />,
	pendingComponent: () => <Page />,
	errorComponent: ({ error }) => <Page error={(error as Error).message} />,
});

function Page({ sample, error }: { sample?: Sample; error?: string }) {
	const router = useRouter();
	return (
		<>
			<title>gcsim - local sample</title>
			<SamplePage
				sample={sample ?? null}
				error={error ?? null}
				retry={() => router.invalidate()}
			/>
		</>
	);
}
