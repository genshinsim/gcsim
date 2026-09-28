import { useNavigate } from "@tanstack/react-router";
import { useCallback } from "react";
import { useDraft } from "../../Stores/AppState";
import type { SendOptions } from "../../Stores/draft";

export function useSendToSimulator() {
	const { send } = useDraft();
	const navigate = useNavigate();
	return useCallback(
		(cfg: string, opts: SendOptions) => {
			send(cfg, opts);
			navigate({ to: "/simulator" });
		},
		[send, navigate],
	);
}
