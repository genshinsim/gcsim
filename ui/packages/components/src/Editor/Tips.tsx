import { Alert, AlertDescription, Button } from "@gcsim/primitives";
import { Trans } from "react-i18next";

export interface TipProps {
	onHide?: () => void;
}

function HideButton({ onHide }: TipProps) {
	if (!onHide) {
		return null;
	}
	return (
		<div className="ml-auto">
			<Button variant="secondary" size="sm" onClick={onHide}>
				<Trans>simple.hide_all_tips</Trans>
			</Button>
		</div>
	);
}

export function TeamTip({ onHide }: TipProps) {
	return (
		<Alert>
			<AlertDescription className="flex flex-col gap-2 text-current">
				<span>
					<Trans>simple.video_pre</Trans>
					<button type="button">
						<Trans>simple.video</Trans>
					</button>
					<Trans>simple.video_post</Trans>
				</span>
				<HideButton onHide={onHide} />
			</AlertDescription>
		</Alert>
	);
}

export function ActionListTip({ onHide }: TipProps) {
	return (
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
				<HideButton onHide={onHide} />
			</AlertDescription>
		</Alert>
	);
}
