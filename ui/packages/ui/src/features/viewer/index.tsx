import { useExecutor, useRunResult, useValidation } from "@gcsim/components";
import type { model } from "@gcsim/types";
import {
	type LinkProps,
	useLocation,
	useNavigate,
	useRouter,
	useSearch,
} from "@tanstack/react-router";
import { usePrefs } from "@ui/stores/AppState";
import axios from "axios";
import { useEffect, useMemo, useRef, useState } from "react";
import { useSendToSimulator } from "../../components/buttons/useSendToSimulator";
import { autoSampleSeed, useSample } from "../sample/useSample";
import { ResultSource } from "./components/LoadingToast";
import type { SignedResult } from "./components/Share";
import type { ViewerTab } from "./search";
import type { ViewerEditorState } from "./tabs/Config";
import UpgradeDialog from "./UpgradeDialog";
import Viewer, { type ViewerActions } from "./Viewer";

type ViewerProps = {
	gitCommit: string;
	mode: string;
};

export type LoadedResult = {
	data: model.SimulationResult;
	signed: SignedResult;
};

export async function loadResult(url: string): Promise<LoadedResult> {
	const resp = await axios.get<string>(url, {
		timeout: 30000,
		responseType: "text",
	});
	return {
		data: JSON.parse(resp.data),
		signed: {
			raw: resp.data,
			hash: resp.headers["x-gcsim-share-auth"] ?? null,
		},
	};
}

function shareResult({ raw, hash }: SignedResult): Promise<string> {
	return axios
		.post("/api/share", raw, {
			headers: {
				"Content-Type": "application/json",
				"X-GCSIM-SHARE-AUTH": hash ?? "",
			},
			transformRequest: (data) => data,
		})
		.then((resp) => link("sh", resp.data));
}

export const WebViewer = ({ mode, gitCommit }: ViewerProps) => {
	const { busy } = useExecutor();
	const { result, raw, hash, config, error } = useRunResult();
	const signed = useMemo(
		() => (raw == null ? null : { raw, hash }),
		[raw, hash],
	);
	return (
		<UpgradableViewer
			data={result}
			signed={signed}
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

export type LoadedViewerProps = ViewerProps & {
	result?: LoadedResult;
	error?: string;
};

export const LoadedViewer = ({
	result,
	error,
	mode,
	gitCommit,
}: LoadedViewerProps) => {
	const router = useRouter();
	return (
		<UpgradableViewer
			data={result?.data ?? null}
			signed={result?.signed ?? null}
			recoveryConfig={null}
			error={error ?? null}
			src={ResultSource.Loaded}
			retry={() => router.invalidate()}
			running={false}
			redirect="/"
			mode={mode}
			gitCommit={gitCommit}
		/>
	);
};

type UpgradableViewerProps = {
	data: model.SimulationResult | null;
	signed: SignedResult | null;
	recoveryConfig: string | null;
	error: string | null;
	src: ResultSource;
	running: boolean;
	redirect: LinkProps["to"];
	mode: string;
	gitCommit: string;
	retry?: () => void;
};

const UpgradableViewer = (props: UpgradableViewerProps) => {
	const { data, running } = props;
	const location = useLocation();
	const navigate = useNavigate();
	const search = useSearch({ strict: false });
	const { run, cancel } = useExecutor();
	const onSendToSimulator = useSendToSimulator();
	const tab = search.tab ?? "results";
	const [linkSeed] = useState(search.seed ?? null);
	const { sampleOnLoad } = usePrefs();
	const sample = useSample({
		config: data?.config_file,
		autoSeed: autoSampleSeed(linkSeed, sampleOnLoad, data?.sample_seed),
	});
	const editor = useViewerEditor(data?.config_file, running);
	useScrollToLocation();

	const actions = useMemo<ViewerActions>(
		() => ({
			onRun: run,
			onSendToSimulator,
			onShare: shareResult,
		}),
		[run, onSendToSimulator],
	);

	const setTab = (next: ViewerTab) =>
		navigate({ to: ".", search: (prev) => ({ ...prev, tab: next }) });

	const generate = (seed: string) => {
		navigate({
			to: ".",
			search: (prev) => ({ ...prev, seed }),
			replace: true,
		});
		sample.generate(seed);
	};

	return (
		<>
			<Viewer
				result={data}
				signed={props.signed}
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
				onClose={() => navigate({ to: props.redirect })}
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
): ViewerEditorState {
	const { isReady } = useExecutor();
	const [config, setConfig] = useState(resultConfig ?? "");
	const { isValid, error } = useValidation(config);

	useEffect(() => {
		setConfig(resultConfig ?? "");
	}, [resultConfig]);

	return {
		config,
		setConfig,
		error,
		canRun: isReady && isValid && !running,
		busy: !isReady || running,
	};
}

function useScrollToLocation() {
	const scrolled = useRef(false);
	const { state, hash } = useLocation();
	const key = state.__TSR_key;
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

function link(route: string, id: string): string {
	return `${window.location.protocol}//${window.location.host}/${route}/${id}`;
}

function extractFromLocation(location: string) {
	if (location.startsWith("/sh/")) {
		return link("sh", location.substring(location.lastIndexOf("/") + 1));
	}
	return null;
}
