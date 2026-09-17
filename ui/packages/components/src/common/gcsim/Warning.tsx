import { dynamicKey } from "@gcsim/localization";
import { Alert, AlertDescription, AlertTitle, Button } from "@gcsim/primitives";
import { Cross2Icon } from "@radix-ui/react-icons";
import React from "react";
import { Trans, useTranslation } from "react-i18next";
import tanuki from "../../images/tanuki.png";
import { cn } from "../../lib/utils";

type WarningVariant = NonNullable<
	React.ComponentProps<typeof Alert>["variant"]
>;

interface WarningProps {
	hideKey: string;
	headerKey: string;
	bodyKey: string;
	bodyComponents?: Record<string, React.ReactElement>;
	showButton?: boolean;
	variant?: WarningVariant;
	className?: string;
}

export function Warning({
	hideKey,
	headerKey,
	bodyKey,
	bodyComponents,
	showButton = true,
	variant = "default",
	className,
}: WarningProps) {
	const { t } = useTranslation();
	const [hide, setHide] = React.useState<boolean>(() => {
		return localStorage.getItem(hideKey) === "true";
	});
	React.useEffect(() => {
		localStorage.setItem(hideKey, hide.toString());
	}, [hide, hideKey]);

	if (hide) {
		if (!showButton) {
			return null;
		}
		return (
			<div className="flex flex-col py-0 min-[1300px]:w-[1100px]">
				<div className="ml-auto">
					<Button size="sm" variant="outline" onClick={() => setHide(false)}>
						{t("db.readme_show")}
					</Button>
				</div>
			</div>
		);
	}

	return (
		<Alert
			variant={variant}
			className={cn(
				"relative flex flex-col items-center gap-2 min-[1300px]:w-[1100px]",
				className,
			)}
		>
			<Button
				size="icon-xs"
				variant="ghost"
				className="absolute top-2 right-2"
				aria-label={t("viewer.close")}
				onClick={() => setHide(true)}
			>
				<Cross2Icon />
			</Button>
			<div className="flex items-center justify-center gap-3 py-2">
				<img src={tanuki} alt="" className="h-10 w-15" />
				<AlertTitle className="line-clamp-none text-center text-xl font-semibold">
					{t(dynamicKey(headerKey))}
				</AlertTitle>
				<img src={tanuki} alt="" className="h-10 w-15" />
			</div>
			<AlertDescription className="pb-3 leading-5">
				<Trans i18nKey={bodyKey as never} components={bodyComponents}>
					<p />
					<p>{{ rerun: t("viewer.rerun") } as never}</p>
					<p className="font-semibold leading-6" />
				</Trans>
			</AlertDescription>
		</Alert>
	);
}

export const RiskWarning = () => (
	<Warning
		hideKey="hide-warning-risk"
		headerKey="warnings.gcsim_risk_title"
		bodyKey="warnings.gcsim_risk_body"
		bodyComponents={{
			b: <b />,
			p: <p />,
			discordlink: (
				// biome-ignore lint/a11y/useAnchorContent: text injected at runtime by <Trans>
				<a
					href="https://discord.com/invite/m7jvjdxx7q"
					target="_blank"
					rel="noreferrer"
				/>
			),
			githublink: (
				// biome-ignore lint/a11y/useAnchorContent: text injected at runtime by <Trans>
				<a
					href="https://github.com/genshinsim/gcsim"
					target="_blank"
					rel="noreferrer"
				/>
			),
			issueslink: (
				// biome-ignore lint/a11y/useAnchorContent: text injected at runtime by <Trans>
				<a
					href="https://github.com/genshinsim/gcsim/issues"
					target="_blank"
					rel="noreferrer"
				/>
			),
		}}
		variant="destructive"
		showButton={false}
	/>
);
