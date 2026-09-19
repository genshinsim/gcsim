import { Button } from "@gcsim/primitives";
import { BookOpen, Play } from "lucide-react";
import type { CSSProperties } from "react";
import { useTranslation } from "react-i18next";
import { Link } from "react-router-dom";
import { HeroArt } from "./HeroArt";
import { getHero } from "./heroImages";
import "./hero.css";

export function Hero() {
	const { t } = useTranslation();
	const hero = getHero();

	return (
		<section
			className="g-hero"
			style={hero.knobs ? (hero.knobs as CSSProperties) : undefined}
		>
			<HeroArt src={hero.src} />
			<div className="g-hero-content relative z-[1] flex max-w-[540px] flex-col">
				<div>
					<h1 className="mb-[18px] font-g-display text-g-hero font-bold text-g-ink">
						{t("dash.hero_title")}
					</h1>
					<p className="max-w-[480px] text-g-lg text-g-ink-dim">
						{t("dash.hero_subtitle")}
					</p>
				</div>
				<div className="mt-auto flex flex-wrap gap-3 pt-10">
					<Button asChild size="lg">
						<Link to="/simulator">
							<Play /> {t("dash.open_simulator")}
						</Link>
					</Button>
					<Button asChild size="lg" variant="outline">
						<a href="https://docs.csim.app">
							<BookOpen /> {t("dash.read_the_docs")}
						</a>
					</Button>
				</div>
			</div>
		</section>
	);
}
