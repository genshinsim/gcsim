import { Button, cn } from "@gcsim/primitives";
import { useTranslation } from "react-i18next";
import { AiFillGithub } from "react-icons/ai";
import { FaDiscord } from "react-icons/fa";

export function CommunityCta({ className }: { className?: string }) {
	const { t } = useTranslation();

	return (
		<div
			className={cn(
				"flex flex-col justify-center rounded-g-lg border border-g-line p-6 shadow-g-card",
				className,
			)}
			style={{
				background:
					"linear-gradient(135deg, var(--g-accent-weak), var(--g-surface))",
			}}
		>
			<h3 className="mb-2 font-g-display text-g-h3 font-semibold text-g-ink">
				{t("dash.community_title")}
			</h3>
			<p className="mb-4 text-g-body text-g-ink-dim">
				{t("dash.community_body")}
			</p>
			<div className="flex flex-wrap gap-g-base-lg">
				<Button asChild>
					<a
						href="https://discord.gg/kthS9wc46f"
						target="_blank"
						rel="noreferrer"
					>
						<FaDiscord /> {t("dash.join_discord")}
					</a>
				</Button>
				<Button asChild variant="outline">
					<a
						href="https://github.com/genshinsim/gcsim"
						target="_blank"
						rel="noreferrer"
					>
						<AiFillGithub /> {t("dash.github")}
					</a>
				</Button>
			</div>
		</div>
	);
}
