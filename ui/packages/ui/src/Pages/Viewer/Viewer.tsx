import { RiskWarning } from "@gcsim/components";
import type { Executor, ExecutorSupplier } from "@gcsim/executors";
import { dynamicKey } from "@gcsim/localization";
import {
	Alert,
	AlertDescription,
	AlertDialog,
	AlertDialogAction,
	AlertDialogCancel,
	AlertDialogContent,
	AlertDialogFooter,
	AlertDialogHeader,
	AlertDialogTitle,
} from "@gcsim/primitives";
import type { SimResults } from "@gcsim/types";
import CopyToClipboard from "@ui/Components/Buttons/CopyToClipboard";
import SendToSimulator from "@ui/Components/Buttons/SendToSimulator";
import { type RootState, useAppSelector } from "@ui/Stores/store";
import queryString from "query-string";
import { useCallback, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import { useHistory } from "react-router";
import type { ResultSource } from ".";
import LoadingToast from "./Components/LoadingToast";
import ViewerNav from "./Components/ViewerNav";
import Warnings from "./Components/Warnings";
import { simResultsToModel } from "./simResultsToModel";
import ConfigUI, { useConfig } from "./Tabs/Config";
import Results from "./Tabs/Results";
import SampleUI, { useSample } from "./Tabs/Sample";

// The viewer only announces intent through these callbacks; the host supplies the effect
// (redux, routing, share API). A button is rendered only when its callback is provided, so a
// read-only host can pass `{}` (or nothing) to get a viewer with no send/rerun/share buttons.
export type ViewerActions = {
	onSendToSimulator?: (cfg: string, opts: { keepTeam: boolean }) => void;
	onRerun?: (cfg: string) => void;
	onShare?: (data: SimResults, hash: string | null) => Promise<string>;
};

type ViewerProps = {
	running: boolean;
	data: SimResults | null;
	hash: string | null;
	recoveryConfig: string | null;
	error: string | null;
	src: ResultSource;
	redirect: string;
	exec: ExecutorSupplier<Executor>;
	retry?: () => void;
	actions?: ViewerActions;
	existingShareLink?: string | null;
};

// The viewer is read-only against the data. Any mutations to the data (resim) must be performed
// above viewer in the hierarchy tree. The viewer can perform whatever additional calculations it
// wants (linreg, stat optimizations, etc) but these computations are *never* stored in the data and
// only exist as long as the page is loaded.
export default ({
	running,
	data,
	hash = "",
	recoveryConfig,
	error,
	src,
	redirect,
	exec,
	retry,
	actions,
	existingShareLink,
}: ViewerProps) => {
	const { t } = useTranslation();
	const parsed = queryString.parse(location.hash);
	const [tabId, setTabId] = useState((parsed.tab as string) ?? "results");

	const cancel = useCallback(() => exec().cancel(), [exec]);
	const sampler = useCallback(
		(cfg: string, seed: string) => exec().sample(cfg, seed),
		[exec],
	);
	const resetTab = useCallback(() => setTabId("results"), []);

	const { sampleOnLoad } = useAppSelector((state: RootState) => {
		return {
			sampleOnLoad: state.app.sampleOnLoad,
		};
	});

	const sample = useSample(running, data, sampleOnLoad, sampler);
	const config = useConfig(data, exec);
	const modelData = useMemo(
		() => (data != null ? simResultsToModel(data) : null),
		[data],
	);
	const names = useMemo(
		() =>
			data?.character_details?.map((c) =>
				t(dynamicKey("game:character_names." + c.name)),
			),
		[data?.character_details, t],
	);

	const tabs: { [k: string]: React.ReactNode } = {
		results: (
			<Results
				data={data}
				modelData={modelData}
				running={running}
				names={names}
			/>
		),
		config: (
			<ConfigUI
				config={config}
				running={running}
				resetTab={resetTab}
				onRerun={actions?.onRerun}
			/>
		),
		analyze: <div></div>,
		sample: (
			<SampleUI
				sampler={sampler}
				data={data}
				sample={sample}
				running={running}
			/>
		),
	};

	return (
		<div className="flex flex-col flex-grow w-full pb-6">
			<div className="flex flex-col px-2 pt-4 empty:pt-0 2xl:mx-auto items-center justify-center">
				<RiskWarning />
			</div>
			<Warnings data={data} />
			<div className="px-2 py-4 w-full 2xl:mx-auto 2xl:container">
				<ViewerNav
					hash={hash}
					tabState={[tabId, setTabId]}
					data={data}
					running={running}
					actions={actions}
					existingShareLink={existingShareLink}
				/>
			</div>
			<div className="basis-full pt-0 mt-0">{tabs[tabId]}</div>
			<LoadingToast
				cancel={cancel}
				running={running}
				src={src}
				error={error}
				current={data?.statistics?.iterations}
				total={data?.simulator_settings?.iterations}
			/>
			<ErrorAlert
				msg={error}
				recoveryConfig={recoveryConfig}
				redirect={redirect}
				retry={retry}
				onSendToSimulator={actions?.onSendToSimulator}
			/>
		</div>
	);
};

const ErrorAlert = ({
	msg,
	recoveryConfig,
	redirect,
	retry,
	onSendToSimulator,
}: {
	msg: string | null;
	recoveryConfig: string | null;
	redirect: string;
	retry?: () => void;
	onSendToSimulator?: ViewerActions["onSendToSimulator"];
}) => {
	const { t } = useTranslation();
	const history = useHistory();

	return (
		<AlertDialog open={msg != null}>
			<AlertDialogContent>
				<AlertDialogHeader>
					<AlertDialogTitle>{t("viewer.error_encountered")}</AlertDialogTitle>
				</AlertDialogHeader>
				<div className="flex flex-col gap-2 mb-1">
					<Alert variant="destructive">
						<AlertDescription>
							<pre className="whitespace-pre-wrap pl-5">{msg}</pre>
						</AlertDescription>
					</Alert>
					{recoveryConfig != null ? (
						<>
							<CopyToClipboard
								config={recoveryConfig}
								className="hidden ml-[7px] sm:flex"
							/>
							<SendToSimulator
								config={recoveryConfig}
								onSendToSimulator={onSendToSimulator}
							/>
						</>
					) : null}
				</div>
				<AlertDialogFooter>
					{retry != null ? (
						<AlertDialogCancel onClick={() => retry()}>
							{t("viewer.retry")}
						</AlertDialogCancel>
					) : null}
					<AlertDialogAction
						variant="destructive"
						onClick={() => history.push(redirect)}
					>
						{t("viewer.return_to_sim")}
					</AlertDialogAction>
				</AlertDialogFooter>
			</AlertDialogContent>
		</AlertDialog>
	);
};
