import { Button } from "@gcsim/primitives";
import { Link } from "@tanstack/react-router";
import { BookOpen, Play } from "lucide-react";
import { useTranslation } from "react-i18next";
import { CommunityCta } from "./CommunityCta";
import "./dashMobile.css";
import { WhatsNew } from "../../components/WhatsNew";
import { usePrefs } from "../../stores/AppState";
import { HeroArt } from "./HeroArt";
import { getHero } from "./heroImages";
import { KQM_DB_URL } from "./kqm";
import { SectionHead } from "./SectionHead";
import { SharedByOthers } from "./SharedByOthers";

export function DashMobile() {
	const { t } = useTranslation();
	const { hero } = usePrefs();

	return (
		<main className="w-full flex-grow bg-g-canvas text-g-ink">
			<div className="g-m-home">
				<div className="g-m-art">
					<HeroArt src={getHero(hero).src} />
				</div>
				<div className="g-m-over">
					<div className="g-m-body mx-auto w-full max-w-[520px] px-4">
						<SectionHead
							className="mb-4"
							title={t("dash.whats_new")}
							subtitle={t("dash.whats_new_sub")}
							linkLabel={t("dash.all_releases")}
							linkHref="https://github.com/genshinsim/gcsim/releases"
						/>
						<WhatsNew className="g-m-fade-card mb-[14px]" />
						<CommunityCta className="mb-g-section" />
						<SectionHead
							className="mb-4"
							title={t("dash.shared_title")}
							subtitle={t("dash.shared_sub")}
							linkLabel={t("dash.view_all")}
							linkHref={KQM_DB_URL}
						/>
						<SharedByOthers className="pb-4" />
					</div>
					<div className="g-m-ctabar">
						<Button asChild className="flex-1">
							<Link to="/simulator">
								<Play /> {t("dash.open_simulator")}
							</Link>
						</Button>
						<Button asChild variant="outline">
							<a href="https://docs.gcsim.app">
								<BookOpen /> {t("dash.docs")}
							</a>
						</Button>
					</div>
				</div>
			</div>
		</main>
	);
}
