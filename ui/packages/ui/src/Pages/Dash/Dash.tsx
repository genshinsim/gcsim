import { DBCard, LatestVersion } from "@gcsim/components";
import { Button, Card } from "@gcsim/primitives";
import type { db } from "@gcsim/types";
import axios from "axios";
import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { Link } from "react-router-dom";

const randQuery = {
	query: {
		$sampleRate: 0.01,
	},
	limit: 1,
	skip: 0,
	sort: {
		create_date: -1,
	},
};

export function Dash() {
	const { t } = useTranslation();

	const [data, setData] = useState<db.Entry[]>([]);
	const [dataIsLoaded, setDataIsLoaded] = useState(false);

	useEffect(() => {
		axios(`/api/db?q=${encodeURIComponent(JSON.stringify(randQuery))}`)
			.then((resp: { data: db.Entries }) => {
				if (resp.data && resp.data.data) {
					setData(resp.data.data);
					setDataIsLoaded(true);
				}
			})
			.catch((err) => console.log(t("viewer.error_encountered") + err.message));
	}, [t]);

	return (
		<main className="w-full flex flex-col items-center flex-grow gap-4 py-4 px-4">
			<Link
				to="/simulator"
				role="button"
				className="bp4-button bp4-intent-success !p-3 !rounded-md"
				tabIndex={0}
			>
				<span className="bp4-button-text text-3xl md:text-4xl lg:text-5xl font-semibold">
					{t("dash.get_started")}
				</span>
			</Link>
			{/* mobile Chrome/Safari needs w-full for padding to work properly, mobile Firefox works fine though... */}
			<div className="flex flex-col gap-4 w-full md:w-fit">
				<Card className="flex flex-col gap-4 items-center">
					<h1 className="text-center text-xl md:text-2xl lg:text-4xl">
						<b>{t("dash.users_submitted")}</b>
					</h1>
					<div>
						{dataIsLoaded
							? data.map((e) => (
									<DBCard
										className="border-0"
										key={e._id}
										entry={e}
										skipTags={-1}
										footer={
											<div className="flex flex-row flex-wrap place-content-end mr-2 gap-4">
												<a
													href={"https://gcsim.app/db/" + e._id}
													target="_blank"
													rel="noopener noreferrer"
												>
													<Button className="bg-blue-600">Show Detail</Button>
												</a>
											</div>
										}
									/>
								))
							: t("sim.loading")}
					</div>
					<Button asChild className="!p-3 !rounded-md">
						<a
							href="https://simpact.app/"
							target="_blank"
							rel="noopener noreferrer"
						>
							<span className="text-xl md:text-2xl font-semibold">
								{t("dash.visit_teams_db")}
							</span>
						</a>
					</Button>
				</Card>
				<LatestVersion />
			</div>
		</main>
	);
}
