import { useCallback } from "react";
import { useNavigate } from "react-router";
import { appActions } from "../../Stores/appSlice";
import { useAppDispatch } from "../../Stores/store";

export function useSendToSimulator() {
	const dispatch = useAppDispatch();
	const navigate = useNavigate();
	return useCallback(
		(cfg: string, { keepTeam }: { keepTeam: boolean }) => {
			dispatch(appActions.setCfg({ cfg, keepTeam }));
			navigate("/simulator");
		},
		[dispatch, navigate],
	);
}
