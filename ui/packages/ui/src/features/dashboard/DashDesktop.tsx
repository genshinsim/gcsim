import { useTranslation } from "react-i18next";
import { WhatsNew } from "../../components/WhatsNew";
import { CommunityCta } from "./CommunityCta";
import { Hero } from "./Hero";
import { SectionHead } from "./SectionHead";

export function DashDesktop() {
	const { t } = useTranslation();

	return (
		<main className="w-full flex-grow bg-g-canvas text-g-ink">
			<div className="mx-auto w-full max-w-[1200px] px-g-page">
				<Hero />

				<SectionHead
					className="mb-[18px]"
					title={t("dash.whats_new")}
					subtitle={t("dash.whats_new_sub")}
					linkLabel={t("dash.all_releases")}
					linkHref="https://github.com/genshinsim/gcsim/releases"
				/>
				<div className="grid grid-cols-1 gap-g-base-lg pb-g-section md:grid-cols-[1.1fr_0.9fr]">
					<WhatsNew />
					<CommunityCta />
				</div>
			</div>
		</main>
	);
}
