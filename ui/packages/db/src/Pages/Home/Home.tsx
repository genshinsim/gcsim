import { Button } from "@gcsim/primitives";
import { WhatsNew } from "@gcsim/ui/src/Pages/Dash/WhatsNew";
import { useTranslation } from "react-i18next";
import { FaCalculator, FaDatabase } from "react-icons/fa";
import { useLocation } from "wouter";

export const Home = () => {
	const { t } = useTranslation();
	const [_, to] = useLocation();

	return (
		<div className="mx-auto flex max-w-[1160px] flex-col gap-g-section p-g-page">
			<section className="flex flex-col gap-g-base-lg rounded-g-xl border border-g-line-soft bg-g-surface p-8">
				<span className="text-g-xs font-semibold uppercase tracking-wide text-g-accent">
					gcsim · community database
				</span>
				<h1 className="font-g-display text-g-h1 font-bold text-g-ink md:text-g-hero">
					{t("db.home.welcome")}
				</h1>
				<p className="max-w-xl text-g-lg text-g-ink-dim">
					{t("db.home.simpact_desc")}
				</p>
				<div className="flex flex-wrap gap-g-base">
					<Button size="lg" onClick={() => to("/database")}>
						<FaDatabase size={14} /> Browse database
					</Button>
					<Button size="lg" variant="ghost" asChild>
						<a
							href="https://gcsim.app/simulator"
							target="_blank"
							rel="noreferrer"
						>
							<FaCalculator size={14} /> Run simulations
						</a>
					</Button>
				</div>
			</section>

			<section className="flex flex-col gap-g-base">
				<h2 className="font-g-display text-g-h2 font-semibold text-g-ink">
					What's new
				</h2>
				<WhatsNew />
			</section>
		</div>
	);
};
