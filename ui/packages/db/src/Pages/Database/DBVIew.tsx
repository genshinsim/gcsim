import { ActionBar } from "SharedComponents/ActionBar";
import { RiskWarning, Warning } from "@gcsim/components";
import type { db } from "@gcsim/types";
import eula from "images/eula.png";
import { useTranslation } from "react-i18next";
import InfiniteScroll from "react-infinite-scroll-component";
import { ListView } from "../../SharedComponents/ListView";

type Props = {
	data: db.Entry[];
	fetchData: () => void;
	hasMore: boolean;
};

export const DBView = (props: Props) => {
	const { t } = useTranslation();
	return (
		<div className="mx-auto flex max-w-[1160px] flex-col gap-g-base-lg px-8 py-4">
			<ActionBar simCount={props.data.length} />
			<RiskWarning />
			<Warning
				hideKey="hide-warning-db"
				headerKey="db.readme_header"
				bodyKey="db.readme_body"
			/>
			{props.data.length === 0 ? (
				<div className="flex h-screen flex-col items-center justify-center">
					<img
						src={eula}
						alt=""
						className="size-32 object-contain opacity-50"
					/>
				</div>
			) : (
				<InfiniteScroll
					dataLength={props.data.length} //This is important field to render the next data
					next={props.fetchData}
					hasMore={props.hasMore}
					loader={<h4>{t("sim.loading")}</h4>}
					endMessage={
						<>
							<p className="mt-4 text-center">
								<b>{t("db.seen_it_all")}</b>
							</p>
							<p className="text-center">{t("db.not_find")}</p>
						</>
					}
					//TODO: enable pull down functionality for refreshing maybe??
				>
					<ListView data={props.data} />
				</InfiniteScroll>
			)}
		</div>
	);
};
