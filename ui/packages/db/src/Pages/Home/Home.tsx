import { LatestVersion } from "@gcsim/components";
import tagData from "@gcsim/data/src/tags.json";
import { dynamicKey } from "@gcsim/localization";
import { Button, Card } from "@gcsim/primitives";
import { useTranslation } from "react-i18next";
import { useLocation } from "wouter";

export const Home = () => {
	const { t } = useTranslation();

	const [_, to] = useLocation();
	const sortedTagnames = Object.keys(tagData)
		.filter((key) => {
			return key !== "0" && key !== "2";
		})
		.map((key) => {
			let name = tagData[key]["display_name"];
			if (key === "1") {
				name = "(Not Tagged)";
			}
			return (
				<li key={key}>
					<span className="font-semibold text-rose-600">{name}</span>
					{`: ${t(dynamicKey("db.home.tag_desc_" + key))}`}
				</li>
			);
		});
	return (
		<main className="w-full flex flex-col items-center flex-grow gap-4 py-4 px-4">
			<Button className="h-auto p-3 rounded-md" onClick={() => to("/database")}>
				<span className="text-3xl md:text-4xl lg:text-5xl font-semibold">
					{t("db.home.get_started")}
				</span>
			</Button>
			<div className="flex flex-col gap-4 w-full md:w-fit">
				<Card className="flex flex-col gap-4 items-center p-4">
					<h1 className="text-center text-xl md:text-2xl lg:text-4xl font-bold">
						{t("db.home.welcome")}
					</h1>
					<div className="min-[1300px]:w-[1225px]">
						<p className="m-2">{t("db.home.simpact_desc")}</p>
						<p className="m-2">{t("db.home.simpact_tag_desc")} </p>
						<p className="m-2">{t("db.home.simpact_tag_list_header")} </p>
						<ul className="list-disc m-4 ml-8">{sortedTagnames}</ul>
					</div>
				</Card>
				<LatestVersion />
			</div>
		</main>
	);
};
