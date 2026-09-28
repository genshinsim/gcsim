import { type EditorProps, ResultsView, RiskWarning } from "@gcsim/components";
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
import type { model } from "@gcsim/types";
import CopyToClipboard from "@ui/Components/Buttons/CopyToClipboard";
import SendToSimulator from "@ui/Components/Buttons/SendToSimulator";
import { useMemo } from "react";
import { useTranslation } from "react-i18next";
import type { SampleState } from "../Sample/useSample";
import LoadingToast, { type ResultSource } from "./Components/LoadingToast";
import ViewerNav from "./Components/ViewerNav";
import Warnings from "./Components/Warnings";
import ConfigUI from "./Tabs/Config";
import SampleUI from "./Tabs/Sample";

export type ViewerActions = {
	onRun?: (cfg: string) => void;
	onSendToSimulator?: (cfg: string, opts: { keepTeam: boolean }) => void;
	onShare?: (
		data: model.SimulationResult,
		hash: string | null,
	) => Promise<string>;
};

type ViewerProps = {
	result: model.SimulationResult | null;
	hash: string | null;
	running: boolean;
	src: ResultSource;
	error: string | null;
	recoveryConfig: string | null;
	tab: string;
	onTabChange: (tab: string) => void;
	editor: Omit<EditorProps, "onRun">;
	sample: SampleState;
	shareLink?: string | null;
	actions?: ViewerActions;
	onCancel: () => void;
	onRetry?: () => void;
	onClose: () => void;
};

// The viewer is read-only against the data. Any mutations to the data (resim) must be performed
// above viewer in the hierarchy tree. The viewer can perform whatever additional calculations it
// wants (linreg, stat optimizations, etc) but these computations are *never* stored in the data and
// only exist as long as the page is loaded.
export default ({
	result,
	hash,
	running,
	src,
	error,
	recoveryConfig,
	tab,
	onTabChange,
	editor,
	sample,
	shareLink,
	actions,
	onCancel,
	onRetry,
	onClose,
}: ViewerProps) => {
	const { t } = useTranslation();
	const names = useMemo(
		() =>
			result?.character_details?.map((c) =>
				t(dynamicKey("game:character_names." + c.name)),
			),
		[result?.character_details, t],
	);

	const onRun = actions?.onRun;
	const tabs: { [k: string]: React.ReactNode } = {
		results: <ResultsView model={result} names={names} />,
		config: (
			<ConfigUI
				{...editor}
				loading={result?.config_file == null}
				canRun={editor.canRun && onRun != null}
				onRun={() => {
					onTabChange("results");
					onRun?.(editor.config);
				}}
			/>
		),
		analyze: <div></div>,
		sample: <SampleUI data={result} sample={sample} running={running} />,
	};

	return (
		<div className="flex flex-col flex-grow w-full pb-6">
			<div className="flex flex-col px-2 pt-4 empty:pt-0 2xl:mx-auto items-center justify-center">
				<RiskWarning />
			</div>
			<Warnings data={result} />
			<div className="px-2 py-4 w-full 2xl:mx-auto 2xl:container">
				<ViewerNav
					hash={hash}
					tabState={[tab, onTabChange]}
					data={result}
					running={running}
					actions={actions}
					existingShareLink={shareLink}
				/>
			</div>
			<div className="basis-full pt-0 mt-0">{tabs[tab]}</div>
			<LoadingToast
				cancel={onCancel}
				running={running}
				src={src}
				error={error}
				current={result?.statistics?.iterations}
				total={result?.simulator_settings?.iterations}
			/>
			<ErrorAlert
				msg={error}
				recoveryConfig={recoveryConfig}
				onRetry={onRetry}
				onClose={onClose}
				onSendToSimulator={actions?.onSendToSimulator}
			/>
		</div>
	);
};

const ErrorAlert = ({
	msg,
	recoveryConfig,
	onRetry,
	onClose,
	onSendToSimulator,
}: {
	msg: string | null;
	recoveryConfig: string | null;
	onRetry?: () => void;
	onClose: () => void;
	onSendToSimulator?: ViewerActions["onSendToSimulator"];
}) => {
	const { t } = useTranslation();

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
					{onRetry != null ? (
						<AlertDialogCancel onClick={() => onRetry()}>
							{t("viewer.retry")}
						</AlertDialogCancel>
					) : null}
					<AlertDialogAction variant="destructive" onClick={onClose}>
						{t("viewer.return_to_sim")}
					</AlertDialogAction>
				</AlertDialogFooter>
			</AlertDialogContent>
		</AlertDialog>
	);
};
