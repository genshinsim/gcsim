import { ButtonGroup, Tabs, TabsList, TabsTrigger } from "@gcsim/primitives";
import type { model } from "@gcsim/types";
import { Link } from "@tanstack/react-router";
import classNames from "classnames";
import { type MouseEvent, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { CopyToClipboard, SendToSimulator } from "../../../Components/Buttons";
import { VIEWER_TABS, type ViewerTab } from "../search";
import type { ViewerActions } from "../Viewer";
import Share from "./Share";

const btnClass = classNames("hidden ml-[7px] sm:flex");

type NavProps = {
	data: model.SimulationResult | null;
	hash: string | null;
	tabState: [ViewerTab, (tab: ViewerTab) => void];
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
		<Tabs value={tabId} onValueChange={(v) => setTabId(v as ViewerTab)}>
			<div className="flex flex-row items-center justify-between gap-2">
				<TabsList variant="line">
					{VIEWER_TABS.map((tab) => (
						<TabsTrigger key={tab} value={tab} asChild>
							<Link
								to="."
								search={(prev) => ({ ...prev, tab })}
								onMouseDown={keepModifierClickInNewTab}
							>
								{t(`viewer.${tab}`)}
							</Link>
						</TabsTrigger>
					))}
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
