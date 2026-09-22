import LatestCharactersData from "@gcsim/data/src/latest_chars.json";
import { dynamicKey } from "@gcsim/localization";
import { Button, Card } from "@gcsim/primitives";
import axios from "axios";
import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import { CharacterTile } from "../../Cards/CharacterTile/CharacterTile";

const majorVersionRegex = /v\d+\.\d+/gm;

export function LatestVersion() {
	const { t } = useTranslation();

	const [{ isLoaded, text, tag, portraits }, setState] = useState({
		isLoaded: false,
		text: "",
		tag: "",
		portraits: [] as string[],
	});

	useEffect(() => {
		axios("https://api.github.com/repos/genshinsim/gcsim/releases/latest")
			.then((resp: { data }) => {
				const majorVersion = majorVersionRegex.exec(resp.data.name);
				let portraits: string[] = [];
				if (majorVersion && majorVersion[0]) {
					portraits = LatestCharactersData[majorVersion[0]] || [];
				}
				setState({
					isLoaded: true,
					text: resp.data.body,
					tag: resp.data.name,
					portraits,
				});
			})
			.catch((err) => console.log(t("viewer.error_encountered") + err.message));
	}, [t]);

	return (
		<Card className="flex flex-col items-center gap-4 overflow-x-auto px-6">
			{isLoaded ? (
				<>
					<div className="flex flex-col gap-4">
						<h1 className="text-center text-g-h1">
							<b>{t("dash.latest_release")}</b>
							<a
								href={`https://github.com/genshinsim/gcsim/releases/tag/${tag}`}
							>
								{tag}
							</a>
						</h1>
					</div>
					<div className="flex flex-col">
						<h2 className="text-center text-g-h2">
							{t("dash.new_characters")}
						</h2>
						<div className="flex gap-4">
							{portraits.map((char) => (
								<div key={char} className="flex flex-col items-center">
									<CharacterTile
										char={{ name: char }}
										i={0}
										invalid={false}
										onImageLoaded={() => {}}
										hideDetails
									/>
									{t(dynamicKey(`game:character_names.${char}`))}
								</div>
							))}
						</div>
					</div>
					<div className="self-start">
						<ReactMarkdown remarkPlugins={[remarkGfm]}>{text}</ReactMarkdown>
					</div>
					<Button asChild size="lg" className="h-auto p-3">
						<a
							href="https://github.com/genshinsim/gcsim/releases"
							target="_blank"
							rel="noreferrer"
						>
							<span className="text-g-lg font-semibold">
								{t("dash.view_releases")}
							</span>
						</a>
					</Button>
				</>
			) : (
				<>{t("sim.loading")}</>
			)}
		</Card>
	);
}
