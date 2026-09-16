import type { IStat } from "@gcsim/types";
import { useTranslation } from "react-i18next";
import { OmniSelect } from "./OmniSelect";
import { statLabel, statPredicate, stats } from "./stats";
import { renderLabeledItem } from "./utils";

export type StatSelectProps = {
	isOpen: boolean;
	onClose: () => void;
	onSelect: (stat: IStat) => void;
	value?: IStat;
};

export function StatSelect({
	isOpen,
	onClose,
	onSelect,
	value,
}: StatSelectProps) {
	const { t } = useTranslation();
	return (
		<OmniSelect<IStat>
			isOpen={isOpen}
			onClose={onClose}
			items={stats}
			itemKey={(stat) => stat}
			itemPredicate={statPredicate}
			itemRenderer={(stat, state) =>
				renderLabeledItem(stat, statLabel(stat), state)
			}
			onSelect={onSelect}
			value={value}
			title={t("simple.stats")}
			placeholder={t("db.type_to_search")}
		/>
	);
}
