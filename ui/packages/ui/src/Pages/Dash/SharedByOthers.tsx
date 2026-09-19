import { cn } from "@gcsim/primitives";
import type { db } from "@gcsim/types";
import axios from "axios";
import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { TeamCard } from "./TeamCard";

const sharedQuery = {
	query: {
		$sampleRate: 0.02,
	},
	limit: 3,
	skip: 0,
	sort: {
		create_date: -1,
	},
};

type SharedByOthersProps = {
	className?: string;
};

export function SharedByOthers({ className }: SharedByOthersProps) {
	const { t } = useTranslation();

	const [entries, setEntries] = useState<db.Entry[]>([]);
	const [isLoaded, setIsLoaded] = useState(false);

	useEffect(() => {
		axios(`/api/db?q=${encodeURIComponent(JSON.stringify(sharedQuery))}`)
			.then((resp: { data: db.Entries }) => {
				if (resp.data?.data) {
					setEntries(resp.data.data);
				}
				setIsLoaded(true);
			})
			.catch((err) => console.log(err));
	}, []);

	if (!isLoaded) {
		return <div className="text-g-sm text-g-ink-mute">{t("sim.loading")}</div>;
	}

	if (entries.length === 0) {
		return null;
	}

	return (
		<div
			className={cn("grid grid-cols-1 gap-g-base-lg md:grid-cols-3", className)}
		>
			{entries.map((e) => (
				<TeamCard entry={e} key={e._id} />
			))}
		</div>
	);
}
