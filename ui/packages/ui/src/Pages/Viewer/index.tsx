import {
	type EditorProps,
	useExecutor,
	useRunResult,
	useValidation,
} from "@gcsim/components";
import type { model } from "@gcsim/types";
import { usePrefs } from "@ui/Stores/AppState";
import axios from "axios";
import queryString from "query-string";
import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { useLocation, useNavigate } from "react-router";
import ExecutorSettingsButton from "../../Components/Buttons/ExecutorSettingsButton";
import { useSendToSimulator } from "../../Components/Buttons/useSendToSimulator";
import { autoSampleSeed, useSample } from "../Sample/useSample";
import { useEditorPrefs } from "../Simulator/editorPrefs";
import { ResultSource } from "./Components/LoadingToast";
import UpgradeDialog from "./UpgradeDialog";
import Viewer, { type ViewerActions } from "./Viewer";

type ViewerProps = {
	gitCommit: string;
	mode: string;
	id?: string; // only used in share
};

export const ShareViewer = (props: ViewerProps) => (
	<FromUrl
		url={"/api/share/" + props.id}
		redirect="/"
		mode={props.mode}
		gitCommit={props.gitCommit}
	/>
);

export const DBViewer = (props: ViewerProps) => (
	<FromUrl
		url={"/api/share/db/" + props.id}
		redirect="/"
		mode={props.mode}
		gitCommit={props.gitCommit}
	/>
);

export const LocalViewer = (props: ViewerProps) => (
	<FromUrl
		url="http://127.0.0.1:8381/data"
		redirect="/"
		mode={props.mode}
		gitCommit={props.gitCommit}
	/>
);

export const WebViewer = ({ mode, gitCommit }: ViewerProps) => {
	const { busy } = useExecutor();
	const { result, hash, config, error } = useRunResult();
	return (
		<UpgradableViewer
			data={result}
			hash={hash}
			recoveryConfig={config}
			error={error}
			src={ResultSource.Generated}
			running={busy}
			redirect="/simulator"
			mode={mode}
			gitCommit={gitCommit}
		/>
	);
};

type FromUrlProps = {
	redirect: string;
	url: string;
	mode: string;
	gitCommit: string;
};

const FromUrl = ({ url, redirect, mode, gitCommit }: FromUrlProps) => {
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
	retry?: () => void;
};

const UpgradableViewer = (props: UpgradableViewerProps) => {
	const { data, running } = props;
	const location = useLocation();
	const navigate = useNavigate();
	const { run, cancel } = useExecutor();
	const onSendToSimulator = useSendToSimulator();
	const [tab, setTab] = useState(
		() => hashParam(location.hash, "tab") ?? "results",
	);
	const [linkSeed] = useState(() => hashParam(location.hash, "sample"));
	const { sampleOnLoad } = usePrefs();
	const sample = useSample({
		config: data?.config_file,
		autoSeed: autoSampleSeed(linkSeed, sampleOnLoad, data?.sample_seed),
		running,
	});
	const editor = useViewerEditor(data?.config_file, running);
	useScrollToLocation();

	const actions = useMemo<ViewerActions>(
		() => ({
			onRun: run,
			onSendToSimulator,
			onShare: (data: model.SimulationResult, hash: string | null) =>
				axios
					.post("/api/share", data, {
						headers: { "X-GCSIM-SHARE-AUTH": hash ?? "" },
					})
					.then((resp) => link("sh", resp.data)),
		}),
		[run, onSendToSimulator],
	);

	const generate = (seed: string) => {
		const parsed = queryString.parse(window.location.hash);
		parsed.sample = seed;
		window.location.hash = queryString.stringify(parsed);
		sample.generate(seed);
	};

	return (
		<>
			<Viewer
				result={data}
				hash={props.hash}
				running={running}
				src={props.src}
				error={props.error}
				recoveryConfig={props.recoveryConfig}
				tab={tab}
				onTabChange={setTab}
				editor={editor}
				sample={{ ...sample, generate }}
				shareLink={extractFromLocation(location.pathname)}
				actions={actions}
				onCancel={cancel}
				onRetry={props.retry}
				onClose={() => navigate(props.redirect)}
			/>
			<UpgradeDialog
				data={data}
				redirect={props.redirect}
				mode={props.mode}
				commit={props.gitCommit}
			/>
		</>
	);
};

function useViewerEditor(
	resultConfig: string | undefined,
	running: boolean,
): Omit<EditorProps, "onRun"> {
	const { isReady } = useExecutor();
	const [config, setConfig] = useState(resultConfig ?? "");
	const [prefs, setPrefs] = useEditorPrefs();
	const { isValid, error, parsedTeam } = useValidation(config);

	useEffect(() => {
		setConfig(resultConfig ?? "");
	}, [resultConfig]);

	return {
		config,
		setConfig,
		error,
		parsedTeam,
		settings: <ExecutorSettingsButton />,
		canRun: isReady && isValid && !running,
		busy: !isReady || running,
		prefs,
		onPrefsChange: setPrefs,
	};
}

function useScrollToLocation() {
	const scrolled = useRef(false);
	const { key, hash } = useLocation();
	const prevKey = useRef(key);

	useEffect(() => {
		if (hash == null) {
			return;
		}

		if (prevKey.current !== key) {
			prevKey.current = key;
			scrolled.current = false;
		}

		if (scrolled.current) {
			return;
		}
		const id = hash.replace("#", "");
		if (!id) {
			return;
		}
		const element = document.getElementById(id);
		if (element) {
			element.scrollIntoView({ behavior: "smooth" });
			scrolled.current = true;
		}
	});
}

function hashParam(hash: string, key: string): string | null {
	const value = queryString.parse(hash)[key];
	return typeof value === "string" ? value : null;
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
