import { cn } from "@gcsim/primitives";
import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { KqmLeaveWarning, useKqmLinkClick } from "./KqmLeaveWarning";
import { KQM_DB_URL, type KqmEntry, parseKqmEntries } from "./kqm";
import { SectionHead } from "./SectionHead";
import { TeamCard } from "./TeamCard";

const sharedQuery = encodeURIComponent(JSON.stringify({ limit: 3 }));

type SharedByOthersProps = {
	headClassName?: string;
	className?: string;
};

export function SharedByOthers({
	headClassName,
	className,
}: SharedByOthersProps) {
	return (
		<KqmLeaveWarning>
			<SharedHead className={headClassName} />
			<SharedEntries className={className} />
		</KqmLeaveWarning>
	);
}

function SharedHead({ className }: { className?: string }) {
	const { t } = useTranslation();
	const onLinkClick = useKqmLinkClick();

	return (
		<SectionHead
			className={className}
			title={t("dash.shared_title")}
			subtitle={t("dash.shared_sub")}
			linkLabel={t("dash.view_all")}
			linkHref={KQM_DB_URL}
			onLinkClick={onLinkClick}
		/>
	);
}

function SharedEntries({ className }: { className?: string }) {
	const { t } = useTranslation();

	const [entries, setEntries] = useState<KqmEntry[]>([]);
	const [isLoaded, setIsLoaded] = useState(false);

	useEffect(() => {
		fetch(`${KQM_DB_URL}/api/db?q=${sharedQuery}`)
			.then((resp) => {
				if (!resp.ok) throw new Error("Could not load simulations");
				return resp.json();
			})
			.then((json) => setEntries(parseKqmEntries(json)))
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
