import { ButtonGroup, Tab, Tabs } from "@blueprintjs/core";
import type { SimResults } from "@gcsim/types";
import classNames from "classnames";
import { type MouseEvent, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import {
	CopyToClipboard,
	SendToSimulator,
	Share,
} from "../../../Components/Buttons";
import type { ViewerActions } from "../Viewer";

const btnClass = classNames("hidden ml-[7px] sm:flex");

type NavProps = {
	data: SimResults | null;
	hash: string | null;
	tabState: [string, (tab: string) => void];
	running: boolean;
	actions?: ViewerActions;
	existingShareLink?: string | null;
};

export default ({
	tabState,
	data,
	hash,
	running,
	actions,
	existingShareLink,
}: NavProps) => {
	const { t } = useTranslation();
	const [tabId, setTabId] = tabState;
	const shareState = useState<string | null>(existingShareLink ?? null);
	const [, setShareLink] = shareState;

	// biome-ignore lint/correctness/useExhaustiveDependencies: data?.config_file re-runs this on a rerun that keeps the same existingShareLink
	useEffect(() => {
		setShareLink(existingShareLink ?? null);
	}, [existingShareLink, setShareLink, data?.config_file]);

	return (
		<Tabs selectedTabId={tabId} onChange={(s) => setTabId(s as string)}>
			<Tab id="results" className="focus:outline-none">
				{/* biome-ignore lint/a11y/useValidAnchor: intentional nav anchor — href drives URL-hash tab deep-linking (reload/bookmark/share to a tab) and native ctrl/cmd-click open-in-new-tab; a <button> would lose both */}
				<a href="#" onClick={ignoreCtrlClick}>
					{t("viewer.results")}
				</a>
			</Tab>
			<Tab id="config" className="focus:outline-none">
				{/* biome-ignore lint/a11y/useValidAnchor: intentional nav anchor — href drives URL-hash tab deep-linking (reload/bookmark/share to a tab) and native ctrl/cmd-click open-in-new-tab; a <button> would lose both */}
				<a href="#tab=config" onClick={ignoreCtrlClick}>
					{t("viewer.config")}
				</a>
			</Tab>
			<Tab id="sample" className="focus:outline-none">
				{/* biome-ignore lint/a11y/useValidAnchor: intentional nav anchor — href drives URL-hash tab deep-linking (reload/bookmark/share to a tab) and native ctrl/cmd-click open-in-new-tab; a <button> would lose both */}
				<a href="#tab=sample" onClick={ignoreCtrlClick}>
					{t("viewer.sample")}
				</a>
			</Tab>
			<Tabs.Expander />
			<ButtonGroup>
				<CopyToClipboard config={data?.config_file} className={btnClass} />
				<SendToSimulator
					config={data?.config_file}
					onSendToSimulator={actions?.onSendToSimulator}
				/>
				<Share
					shareState={shareState}
					data={data}
					hash={hash}
					running={running}
					className={btnClass}
					onShare={actions?.onShare}
				/>
			</ButtonGroup>
		</Tabs>
	);
};

function ignoreCtrlClick(e: MouseEvent) {
	if (e.ctrlKey) {
		e.stopPropagation();
	}
}
