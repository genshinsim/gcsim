import { ButtonGroup, Position, Tab, Tabs, Toaster } from "@blueprintjs/core";
import type { SimResults } from "@gcsim/types";
import classNames from "classnames";
import { type MouseEvent, useRef, useState } from "react";
import { useTranslation } from "react-i18next";
import {
	CopyToClipboard,
	SendToSimulator,
	Share,
} from "../../../Components/Buttons";

const btnClass = classNames("hidden ml-[7px] sm:flex");

type NavProps = {
	data: SimResults | null;
	hash: string | null;
	tabState: [string, (tab: string) => void];
	running: boolean;
};

export default ({ tabState, data, hash, running }: NavProps) => {
	const { t } = useTranslation();
	const [tabId, setTabId] = tabState;
	const copyToast = useRef<Toaster>(null);
	const shareState = useState<string | null>(null);

	return (
		<Tabs selectedTabId={tabId} onChange={(s) => setTabId(s as string)}>
			<Tab id="results" className="focus:outline-none">
				{/* biome-ignore lint/a11y/useValidAnchor: intentional nav anchor — href drives URL-hash tab deep-linking (reload/bookmark/share to a tab) and native ctrl/cmd-click open-in-new-tab; a <button> would lose both */}
				<a href="#" onClick={ignoreCtrlClick}>
					{t<string>("viewer.results")}
				</a>
			</Tab>
			<Tab id="config" className="focus:outline-none">
				{/* biome-ignore lint/a11y/useValidAnchor: intentional nav anchor — href drives URL-hash tab deep-linking (reload/bookmark/share to a tab) and native ctrl/cmd-click open-in-new-tab; a <button> would lose both */}
				<a href="#tab=config" onClick={ignoreCtrlClick}>
					{t<string>("viewer.config")}
				</a>
			</Tab>
			<Tab id="sample" className="focus:outline-none">
				{/* biome-ignore lint/a11y/useValidAnchor: intentional nav anchor — href drives URL-hash tab deep-linking (reload/bookmark/share to a tab) and native ctrl/cmd-click open-in-new-tab; a <button> would lose both */}
				<a href="#tab=sample" onClick={ignoreCtrlClick}>
					{t<string>("viewer.sample")}
				</a>
			</Tab>
			<Tabs.Expander />
			<ButtonGroup>
				<CopyToClipboard
					copyToast={copyToast}
					config={data?.config_file}
					className={btnClass}
				/>
				<SendToSimulator config={data?.config_file} />
				<Share
					copyToast={copyToast}
					shareState={shareState}
					data={data}
					hash={hash}
					running={running}
					className={btnClass}
				/>
			</ButtonGroup>
			<Toaster ref={copyToast} position={Position.TOP_RIGHT} />
		</Tabs>
	);
};

function ignoreCtrlClick(e: MouseEvent) {
	if (e.ctrlKey) {
		e.stopPropagation();
	}
}
