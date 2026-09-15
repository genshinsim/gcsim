import { useCallback } from "react";
import { useHistory } from "react-router";
import { appActions } from "../../Stores/appSlice";
import { useAppDispatch } from "../../Stores/store";

export function useSendToSimulator() {
	const dispatch = useAppDispatch();
	const history = useHistory();
	return useCallback(
		(cfg: string, { keepTeam }: { keepTeam: boolean }) => {
			dispatch(appActions.setCfg({ cfg, keepTeam }));
			history.push("/simulator");
		},
		[dispatch, history],
	);
}
