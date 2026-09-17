import { Alert, AlertDescription, Button } from "@gcsim/primitives";
import React from "react";
import { Trans } from "react-i18next";
import { useAppDispatch, useAppSelector } from "../../Stores/store";
import { userActions } from "../../Stores/userSlice";

export const ActionListTooltip = () => {
	const settings = useAppSelector((state) => state.user.data.settings);
	const dispatch = useAppDispatch();
	const toggleTips = () => {
		dispatch(
			userActions.setUserSettings({
				showTips: !settings.showTips,
				showBuilder: settings.showBuilder,
				showNameSearch: settings.showNameSearch,
			}),
		);
	};

	if (!settings.showTips) {
		return null;
	}

	return (
		<div className="pl-2 pr-2 pt-2">
			<Alert>
				<AlertDescription className="flex flex-col gap-2 text-current">
					<p>
						<Trans>simple.discord_pre</Trans>
						<a
							href="https://discord.gg/W36ZwwhEaG"
							target="_blank"
							rel="noreferrer"
						>
							Discord
						</a>
						<Trans>simple.discord_post</Trans>
					</p>
					<p>
						<Trans>simple.documentation_pre</Trans>
						<a
							href="https://docs.gcsim.app/guides"
							target="_blank"
							rel="noreferrer"
						>
							<Trans>simple.documentation</Trans>
						</a>
						<Trans>simple.documentation_post</Trans>
					</p>
					<div className="ml-auto">
						<Button variant="secondary" size="sm" onClick={toggleTips}>
							<Trans>simple.hide_all_tips</Trans>
						</Button>
					</div>
				</AlertDescription>
			</Alert>
		</div>
	);
};

export const TeamBuilderTooltip = () => {
	const [, setOpenAddCharHelp] = React.useState<boolean>(false);

	const settings = useAppSelector((state) => state.user.data.settings);
	const dispatch = useAppDispatch();
	const toggleTips = () => {
		dispatch(
			userActions.setUserSettings({
				showTips: !settings.showTips,
				showBuilder: settings.showBuilder,
				showNameSearch: settings.showNameSearch,
			}),
		);
	};

	if (!settings.showTips) {
		return null;
	}

	return (
		<div className="pl-2 pr-2">
			<Alert>
				<AlertDescription className="flex flex-col gap-2 text-current">
					<span>
						<Trans>simple.video_pre</Trans>
						<button type="button" onClick={() => setOpenAddCharHelp(true)}>
							<Trans>simple.video</Trans>
						</button>
						<Trans>simple.video_post</Trans>
					</span>
					<div className="ml-auto">
						<Button variant="secondary" size="sm" onClick={toggleTips}>
							<Trans>simple.hide_all_tips</Trans>
						</Button>
					</div>
				</AlertDescription>
			</Alert>
		</div>
	);
};
