import { cn } from "@gcsim/primitives";
import { db } from "@gcsim/types";
import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { KQM_DB_URL } from "./kqm";
import { TeamCard } from "./TeamCard";

const sharedQuery = encodeURIComponent(JSON.stringify({ limit: 3 }));

type SharedByOthersProps = {
	className?: string;
};

export function SharedByOthers({ className }: SharedByOthersProps) {
	const { t } = useTranslation();

	const [entries, setEntries] = useState<db.Entry[]>([]);
	const [isLoaded, setIsLoaded] = useState(false);

	useEffect(() => {
		fetch(`${KQM_DB_URL}/api/db?q=${sharedQuery}`)
			.then((resp) => {
				if (!resp.ok) throw new Error("Could not load simulations");
				return resp.json();
			})
			.then((json) => setEntries(db.Entries.fromJSON(json).data ?? []))
			.catch((err) => console.log(err))
			.finally(() => setIsLoaded(true));
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
