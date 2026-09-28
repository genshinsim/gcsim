import { useCallback } from "react";
import { useNavigate } from "react-router";
import { draft } from "../../Stores";

export function useSendToSimulator() {
	const navigate = useNavigate();
	return useCallback(
		(cfg: string, opts: { keepTeam: boolean }) => {
			draft.send(cfg, opts);
			navigate("/simulator");
		},
		[navigate],
	);
}
