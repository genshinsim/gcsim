import { Button } from "@gcsim/primitives";
import { WhatsNew } from "@gcsim/ui/src/Pages/Dash/WhatsNew";
import { Trans, useTranslation } from "react-i18next";
import { FaCalculator, FaDatabase } from "react-icons/fa";
import { useLocation } from "wouter";

export const Home = () => {
	const { t } = useTranslation();
	const [_, to] = useLocation();

	return (
		<div className="mx-auto flex max-w-[1160px] flex-col gap-g-section p-g-page">
			<section className="flex flex-col gap-g-base-lg rounded-g-xl border border-g-line-soft bg-g-surface p-8">
				<span className="text-g-xs font-semibold uppercase tracking-wide text-g-accent">
					{t("db.home.eyebrow")}
				</span>
				<h1 className="font-g-display text-g-h1 font-bold text-g-ink md:text-g-hero">
					{t("db.home.welcome")}
				</h1>
				<p className="max-w-xl text-balance text-g-lg text-g-ink-dim">
					{t("db.home.simpact_desc")}
				</p>
				<blockquote className="max-w-xl border-l-2 border-g-accent/60 pl-g-base text-g-sm leading-5 text-g-ink-dim [&>p]:mb-g-base-sm [&>p:last-child]:mb-0">
					<span className="mb-g-base-sm block font-semibold text-g-ink">
						{t("db.readme_header")}
					</span>
					<Trans i18nKey={"db.readme_body" as never}>
						<p />
						<p>{{ rerun: t("viewer.rerun") } as never}</p>
						<p />
					</Trans>
				</blockquote>
				<div className="flex flex-wrap gap-g-base">
					<Button size="lg" onClick={() => to("/database")}>
						<FaDatabase size={14} /> {t("db.home.browse_database")}
					</Button>
					<Button size="lg" variant="ghost" asChild>
						<a
							href="https://gcsim.app/simulator"
							target="_blank"
							rel="noreferrer"
						>
							<FaCalculator size={14} /> {t("db.home.run_simulations")}
						</a>
					</Button>
				</div>
			</section>

			<section className="flex flex-col gap-g-base">
				<h2 className="font-g-display text-g-h2 font-semibold text-g-ink">
					{t("dash.whats_new")}
				</h2>
				<WhatsNew />
			</section>
		</div>
	);
};
