import { CharacterTile } from "@gcsim/components";
import LatestCharactersData from "@gcsim/data/src/latest_chars.json";
import { cn } from "@gcsim/primitives";
import axios from "axios";
import { ArrowRight } from "lucide-react";
import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";

const majorVersionRegex = /v\d+\.\d+/;

type GithubRelease = {
	name: string;
	body: string;
	published_at: string;
	html_url: string;
};

export function WhatsNew({ className }: { className?: string }) {
	const { t } = useTranslation();

	const [{ isLoaded, name, body, publishedAt, htmlUrl, portraits }, setState] =
		useState({
			isLoaded: false,
			name: "",
			body: "",
			publishedAt: "",
			htmlUrl: "",
			portraits: [] as string[],
		});

	useEffect(() => {
		axios("https://api.github.com/repos/genshinsim/gcsim/releases/latest")
			.then((resp: { data: GithubRelease }) => {
				const majorVersion = majorVersionRegex.exec(resp.data.name);
				let portraits: string[] = [];
				if (majorVersion?.[0]) {
					portraits = LatestCharactersData[majorVersion[0]] || [];
				}
				setState({
					isLoaded: true,
					name: resp.data.name,
					body: resp.data.body,
					publishedAt: resp.data.published_at,
					htmlUrl: resp.data.html_url,
					portraits,
				});
			})
			.catch((err) => console.log(err.message));
	}, []);

	const cardClass = cn(
		"border border-g-line bg-g-surface rounded-g-lg p-g-card shadow-g-card",
		className,
	);

	if (!isLoaded) {
		return <div className={cardClass}>{t("sim.loading")}</div>;
	}

	const formattedDate = new Date(publishedAt).toLocaleDateString();

	return (
		<div className={cardClass}>
			<div className="mb-3 flex flex-wrap items-center gap-g-base">
				<span className="inline-flex items-center gap-g-base-sm rounded-g-pill bg-g-accent-weak px-2.5 py-1 font-g-mono text-g-xs font-semibold text-g-accent">
					{name}
				</span>
				<span className="text-g-sm text-g-ink-mute">
					{t("dash.latest_release_label")} · {formattedDate}
				</span>
			</div>

			{portraits.length > 0 && (
				<div className="mb-3 flex flex-wrap gap-g-base">
					{portraits.map((char) => (
						<CharacterTile
							key={char}
							char={{ name: char }}
							i={0}
							invalid={false}
							onImageLoaded={() => {}}
							hideDetails
						/>
					))}
				</div>
			)}

			<div className="max-h-[240px] overflow-y-auto text-g-body text-g-ink-dim leading-relaxed">
				<ReactMarkdown
					remarkPlugins={[remarkGfm]}
					components={{
						a: (props) => (
							<a
								{...props}
								className="text-g-accent hover:text-g-accent-hover"
								target="_blank"
								rel="noreferrer"
							/>
						),
						h1: (props) => (
							<h1
								{...props}
								className="mb-1 mt-2 font-g-display font-semibold text-g-ink"
							/>
						),
						h2: (props) => (
							<h2
								{...props}
								className="mb-1 mt-2 font-g-display font-semibold text-g-ink"
							/>
						),
						h3: (props) => (
							<h3
								{...props}
								className="mb-1 mt-2 font-g-display font-semibold text-g-ink"
							/>
						),
						ul: (props) => <ul {...props} className="list-disc pl-5" />,
						li: (props) => <li {...props} className="text-g-ink-dim" />,
						code: (props) => (
							<code {...props} className="font-g-mono text-g-sm" />
						),
						p: (props) => <p {...props} className="mb-2" />,
					}}
				>
					{body}
				</ReactMarkdown>
			</div>

			<a
				href={htmlUrl}
				target="_blank"
				rel="noreferrer"
				className="mt-3 inline-flex items-center gap-g-base-sm text-g-sm font-semibold text-g-accent hover:text-g-accent-hover"
			>
				{t("dash.full_changelog")} <ArrowRight className="size-3.5" />
			</a>
		</div>
	);
}
