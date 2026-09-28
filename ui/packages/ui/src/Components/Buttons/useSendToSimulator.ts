import { useCallback } from "react";
import { useNavigate } from "react-router";
import { useDraft } from "../../Stores/AppState";
import type { SendOptions } from "../../Stores/draft";

export function useSendToSimulator() {
	const { send } = useDraft();
	const navigate = useNavigate();
	return useCallback(
		(cfg: string, opts: SendOptions) => {
			send(cfg, opts);
			navigate("/simulator");
		},
		[send, navigate],
	);
}
