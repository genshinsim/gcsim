import type { Executor } from "@gcsim/executors";
import type { SimResults } from "@gcsim/types";
import { throttle } from "lodash-es";
import type { AppThunk } from "../../Stores/store";
import { viewerActions } from "../../Stores/viewerSlice";
import { VIEWER_THROTTLE } from "../Viewer";

export function runSim(pool: Executor, cfg: string): AppThunk {
	return (dispatch) => {
		dispatch(viewerActions.start());

		const updateResult = throttle(
			(res: SimResults, hash: string) => {
				dispatch(viewerActions.setResult({ data: res, hash: hash }));
			},
			VIEWER_THROTTLE,
			{ leading: true, trailing: true },
		);

		pool.run(cfg, updateResult).catch((err) => {
			dispatch(viewerActions.setError({ recoveryConfig: cfg, error: err }));
		});
	};
}
