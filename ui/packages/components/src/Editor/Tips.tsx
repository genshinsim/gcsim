import { Alert, AlertDescription, Button } from "@gcsim/primitives";
import { Trans } from "react-i18next";

export interface TipsProps {
	onHide?: () => void;
}

export function Tips({ onHide }: TipsProps) {
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
				{onHide ? (
					<div className="ml-auto">
						<Button variant="secondary" size="sm" onClick={onHide}>
							<Trans>simple.hide_all_tips</Trans>
						</Button>
					</div>
				) : null}
			</AlertDescription>
		</Alert>
	);
}
