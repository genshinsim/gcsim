import { Trans, useTranslation } from "react-i18next";

const linkClass = "font-semibold text-g-accent underline hover:text-g-ink";

export const Home = () => {
	const { t } = useTranslation();

	return (
		<div className="mx-auto flex max-w-[1160px] flex-col p-g-page">
			<section className="flex flex-col gap-g-base-lg rounded-g-xl border border-g-line-soft bg-g-surface p-8">
				<h1 className="font-g-display text-g-h1 font-bold text-g-ink md:text-g-hero">
					{t("db.home.archived_title")}
				</h1>
				<div className="flex flex-col gap-g-base text-g-lg text-g-ink-dim">
					<p>
						<Trans
							i18nKey="db.home.archived_body"
							components={{
								// biome-ignore lint/a11y/useAnchorContent: text injected at runtime by <Trans>
								url: <a href="https://db.kqm.gg/" className={linkClass} />,
							}}
						/>
					</p>
					<p>
						<Trans
							i18nKey="db.home.archived_contact"
							components={{
								url: (
									// biome-ignore lint/a11y/useAnchorContent: text injected at runtime by <Trans>
									<a
										href="https://discord.gg/keqing"
										target="_blank"
										rel="noreferrer"
										className={linkClass}
									/>
								),
							}}
						/>
					</p>
				</div>
			</section>
		</div>
	);
};
