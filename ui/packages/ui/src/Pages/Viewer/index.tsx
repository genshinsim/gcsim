import { useExecutor, useRunResult } from "@gcsim/components";
import type { Executor, ExecutorSupplier } from "@gcsim/executors";
import type { model } from "@gcsim/types";
import axios from "axios";
import { useCallback, useEffect, useMemo, useState } from "react";
import { useLocation } from "react-router";
import { useSendToSimulator } from "../../Components/Buttons/useSendToSimulator";
import UpgradeDialog from "./UpgradeDialog";
import Viewer, { type ViewerActions } from "./Viewer";

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

export const WebViewer = ({ exec, mode, gitCommit }: ViewerProps) => {
	const { busy } = useExecutor();
	const { result, hash, config, error } = useRunResult();
	return (
		<UpgradableViewer
			data={result}
			hash={hash}
			recoveryConfig={config}
			error={error}
			src={ResultSource.Generated}
			exec={exec}
			running={busy}
			redirect="/simulator"
			mode={mode}
			gitCommit={gitCommit}
		/>
	);
};

type FromUrlProps = {
	exec: ExecutorSupplier<Executor>;
	redirect: string;
	url: string;
	mode: string;
	gitCommit: string;
};

const FromUrl = ({ exec, url, redirect, mode, gitCommit }: FromUrlProps) => {
	const [data, setData] = useState<model.SimulationResult | null>(null);
	const [hash, setHash] = useState<string | null>(null);
	const [error, setError] = useState<string | null>(null);

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

	return (
		<UpgradableViewer
			data={data}
			hash={hash}
			recoveryConfig={null}
			error={error}
			src={ResultSource.Loaded}
			exec={exec}
			retry={request}
			running={false}
			redirect={redirect}
			mode={mode}
			gitCommit={gitCommit}
		/>
	);
};

type UpgradableViewerProps = {
	data: model.SimulationResult | null;
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
};

const UpgradableViewer = (props: UpgradableViewerProps) => {
	const actions = useViewerActions();
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
				data={props.data}
				redirect={props.redirect}
				mode={props.mode}
				commit={props.gitCommit}
			/>
		</>
	);
};

function useViewerActions(): ViewerActions {
	const { run } = useExecutor();
	const onSendToSimulator = useSendToSimulator();
	return useMemo(
		() => ({
			onSendToSimulator,
			onRerun: run,
			onShare: (data: model.SimulationResult, hash: string | null) =>
				axios
					.post("/api/share", data, {
						headers: { "X-GCSIM-SHARE-AUTH": hash ?? "" },
					})
					.then((resp) => link("sh", resp.data)),
		}),
		[run, onSendToSimulator],
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
