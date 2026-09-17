import { ButtonGroup, Tabs, TabsList, TabsTrigger } from "@gcsim/primitives";
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
		<Tabs value={tabId} onValueChange={setTabId}>
			<div className="flex flex-row items-center justify-between gap-2">
				<TabsList variant="line">
					<TabsTrigger value="results" asChild>
						{/* biome-ignore lint/a11y/useValidAnchor: intentional nav anchor — href drives URL-hash tab deep-linking (reload/bookmark/share to a tab) and native ctrl/cmd-click open-in-new-tab; a <button> would lose both */}
						<a href="#" onMouseDown={keepModifierClickInNewTab}>
							{t("viewer.results")}
						</a>
					</TabsTrigger>
					<TabsTrigger value="config" asChild>
						{/* biome-ignore lint/a11y/useValidAnchor: intentional nav anchor — href drives URL-hash tab deep-linking (reload/bookmark/share to a tab) and native ctrl/cmd-click open-in-new-tab; a <button> would lose both */}
						<a href="#tab=config" onMouseDown={keepModifierClickInNewTab}>
							{t("viewer.config")}
						</a>
					</TabsTrigger>
					<TabsTrigger value="sample" asChild>
						{/* biome-ignore lint/a11y/useValidAnchor: intentional nav anchor — href drives URL-hash tab deep-linking (reload/bookmark/share to a tab) and native ctrl/cmd-click open-in-new-tab; a <button> would lose both */}
						<a href="#tab=sample" onMouseDown={keepModifierClickInNewTab}>
							{t("viewer.sample")}
						</a>
					</TabsTrigger>
				</TabsList>
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
			</div>
		</Tabs>
	);
};

// A ctrl/cmd-click on a tab must open that tab in a new browser tab (the
// anchor's native behavior) without switching the current tab. Preventing the
// mousedown default both blocks focus and trips Radix's composeEventHandlers
// gate so the trigger skips activation, while the click's new-tab navigation
// still fires.
function keepModifierClickInNewTab(e: MouseEvent) {
	if (e.ctrlKey || e.metaKey) {
		e.preventDefault();
	}
}
