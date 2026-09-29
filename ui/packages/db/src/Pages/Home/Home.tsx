import type { ReactNode } from "react";
import { Trans, useTranslation } from "react-i18next";

const TextLink = ({
	href,
	newTab,
	children,
}: {
	href: string;
	newTab?: boolean;
	children?: ReactNode;
}) => (
	<a
		href={href}
		target={newTab ? "_blank" : undefined}
		rel={newTab ? "noreferrer" : undefined}
		className="font-semibold text-g-accent underline hover:text-g-ink"
	>
		{children}
	</a>
);

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
							components={{ url: <TextLink href="https://db.kqm.gg/" /> }}
						/>
					</p>
					<p>
						<Trans
							i18nKey="db.home.archived_contact"
							components={{
								url: <TextLink href="https://discord.gg/keqing" newTab />,
							}}
						/>
					</p>
				</div>
			</section>
		</div>
	);
};
