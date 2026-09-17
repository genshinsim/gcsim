import type { Executor, ExecutorSupplier } from "@gcsim/executors";
import type { SimResults } from "@gcsim/types";
import axios from "axios";
import { throttle } from "lodash-es";
import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { useLocation, useNavigate } from "react-router";
import { useSendToSimulator } from "../../Components/Buttons/useSendToSimulator";
import {
	type RootState,
	useAppDispatch,
	useAppSelector,
} from "../../Stores/store";
import { viewerActions } from "../../Stores/viewerSlice";
import { runSim } from "../Simulator/Toolbox";
import UpgradeDialog from "./UpgradeDialog";
import Viewer, { type ViewerActions } from "./Viewer";

// TODO: make this flush rate configurable?
export const VIEWER_THROTTLE = 100;

export enum ResultSource {
	Loaded,
	Generated,
}

type ViewerProps = {
	exec: ExecutorSupplier<Executor>;
	gitCommit: string;
	mode: string;
	id?: string; // only used in share
};

export const ShareViewer = (props: ViewerProps) => (
	<FromUrl
		exec={props.exec}
		url={"/api/share/" + props.id}
		redirect="/"
		mode={props.mode}
		gitCommit={props.gitCommit}
	/>
);

export const DBViewer = (props: ViewerProps) => (
	<FromUrl
		exec={props.exec}
		url={"/api/share/db/" + props.id}
		redirect="/"
		mode={props.mode}
		gitCommit={props.gitCommit}
	/>
);

export const LocalViewer = (props: ViewerProps) => (
	<FromUrl
		exec={props.exec}
		url="http://127.0.0.1:8381/data"
		redirect="/"
		mode={props.mode}
		gitCommit={props.gitCommit}
	/>
);

export const WebViewer = (props: ViewerProps) => (
	<FromState
		exec={props.exec}
		redirect="/simulator"
		mode={props.mode}
		gitCommit={props.gitCommit}
	/>
);

function useRunningState(exec: ExecutorSupplier<Executor>): boolean {
	const [isRunning, setRunning] = useState(true);

	useEffect(() => {
		const check = setInterval(() => {
			setRunning(exec().running());
		}, VIEWER_THROTTLE - 50);
		return () => clearInterval(check);
	}, [exec]);

	return isRunning;
}

type FromUrlProps = {
	exec: ExecutorSupplier<Executor>;
	redirect: string;
	url: string;
	mode: string;
	gitCommit: string;
};

const FromUrl = ({ exec, url, redirect, mode, gitCommit }: FromUrlProps) => {
	const [data, setData] = useState<SimResults | null>(null);
	const [hash, setHash] = useState<string | null>(null);
	const [error, setError] = useState<string | null>(null);
	const [src, setSrc] = useState<ResultSource>(ResultSource.Loaded);
	const isRunning = useRunningState(exec);

	const request = useCallback(() => {
		setError(null);
		axios
			.get(url, { timeout: 30000 })
			.then((resp) => {
				setData(resp.data);
				console.log(resp.data);
				setHash(resp.headers["x-gcsim-share-auth"] ?? null);
			})
			.catch((e) => {
				setError(e.message);
			});
	}, [url]);
	useEffect(() => request(), [request]);

	const updateResult = useRef(
		throttle(
			(res: SimResults | null) => {
				setData(res);
				setSrc(ResultSource.Generated);
			},
			VIEWER_THROTTLE,
			{ leading: true, trailing: true },
		),
	);
	const handleSetError = (_: string | null, error: string | null) => {
		if (error == null) {
			return;
		}
		setError(error);
	};

	return (
		<UpgradableViewer
			data={data}
			hash={hash}
			recoveryConfig={null}
			error={error}
			src={src}
			exec={exec}
			retry={request}
			running={isRunning}
			redirect={redirect}
			mode={mode}
			gitCommit={gitCommit}
			setResult={updateResult.current}
			setError={handleSetError}
		/>
	);
};

type FromStateProps = {
	exec: ExecutorSupplier<Executor>;
	redirect: string;
	mode: string;
	gitCommit: string;
};

const FromState = ({ exec, redirect, mode, gitCommit }: FromStateProps) => {
	const isRunning = useRunningState(exec);
	const { data, hash, recoveryConfig, error } = useAppSelector(
		(state: RootState) => {
			return {
				data: state.viewer.data,
				hash: state.viewer.hash,
				recoveryConfig: state.viewer.recoveryConfig,
				error: state.viewer.error,
			};
		},
	);
	const dispatch = useAppDispatch();

	const setResult = useRef(
		throttle(
			(result: SimResults | null, hash: string | null) => {
				if (result == null) {
					return;
				}
				dispatch(viewerActions.setResult({ data: result, hash: hash }));
			},
			VIEWER_THROTTLE,
			{ leading: true, trailing: true },
		),
	);

	const setError = (recoveryConfig: string | null, error: string | null) => {
		if (error == null) {
			return;
		}
		dispatch(viewerActions.setError({ recoveryConfig, error }));
	};

	return (
		<UpgradableViewer
			data={data}
			hash={hash}
			recoveryConfig={recoveryConfig}
			error={error}
			src={ResultSource.Generated}
			exec={exec}
			running={isRunning}
			redirect={redirect}
			mode={mode}
			gitCommit={gitCommit}
			setResult={setResult.current}
			setError={setError}
		/>
	);
};

type UpgradableViewerProps = {
	data: SimResults | null;
	hash: string | null;
	recoveryConfig: string | null;
	error: string | null;
	src: ResultSource;
	running: boolean;
	redirect: string;
	mode: string;
	gitCommit: string;
	exec: ExecutorSupplier<Executor>;
	retry?: () => void;
	setResult: (r: SimResults | null, hash: string | null) => void;
	setError: (recoveryConfig: string | null, err: string | null) => void;
};

const UpgradableViewer = (props: UpgradableViewerProps) => {
	const actions = useViewerActions(props.exec);
	const location = useLocation();
	return (
		<>
			<Viewer
				running={props.running}
				data={props.data}
				hash={props.hash}
				src={props.src}
				recoveryConfig={props.recoveryConfig}
				error={props.error}
				redirect={props.redirect}
				exec={props.exec}
				retry={props.retry}
				actions={actions}
				existingShareLink={extractFromLocation(location.pathname)}
			/>
			<UpgradeDialog
				exec={props.exec}
				data={props.data}
				redirect={props.redirect}
				mode={props.mode}
				commit={props.gitCommit}
				setResult={props.setResult}
				setError={props.setError}
			/>
		</>
	);
};

function useViewerActions(exec: ExecutorSupplier<Executor>): ViewerActions {
	const dispatch = useAppDispatch();
	const navigate = useNavigate();
	const onSendToSimulator = useSendToSimulator();
	return useMemo(
		() => ({
			onSendToSimulator,
			onRerun: (cfg: string) => {
				dispatch(runSim(exec(), cfg));
				navigate("/web");
			},
			onShare: (data: SimResults, hash: string | null) =>
				axios
					.post("/api/share", data, {
						headers: { "X-GCSIM-SHARE-AUTH": hash ?? "" },
					})
					.then((resp) => link("sh", resp.data)),
		}),
		[dispatch, navigate, exec, onSendToSimulator],
	);
}

function link(route: string, id: string): string {
	return `${window.location.protocol}//${window.location.host}/${route}/${id}`;
}

function extractFromLocation(location: string) {
	if (location.startsWith("/sh/")) {
		return link("sh", location.substring(location.lastIndexOf("/") + 1));
	} else if (location.startsWith("/db/")) {
		return link("db", location.substring(location.lastIndexOf("/") + 1));
	}
	return null;
}
