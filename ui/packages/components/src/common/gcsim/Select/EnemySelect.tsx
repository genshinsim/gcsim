import type { IEnemy } from "@gcsim/types";
import { useTranslation } from "react-i18next";
import { enemies, enemyLabel, enemyPredicate } from "./enemies";
import { OmniSelect } from "./OmniSelect";
import { renderLabeledItem } from "./utils";

export type EnemySelectProps = {
	isOpen: boolean;
	onClose: () => void;
	onSelect: (enemy: IEnemy) => void;
	value?: IEnemy;
};

export function EnemySelect({
	isOpen,
	onClose,
	onSelect,
	value,
}: EnemySelectProps) {
	const { t } = useTranslation();
	return (
		<OmniSelect<IEnemy>
			isOpen={isOpen}
			onClose={onClose}
			items={enemies}
			itemKey={(enemy) => enemy}
			itemPredicate={enemyPredicate}
			itemRenderer={(enemy, state) =>
				renderLabeledItem(enemy, enemyLabel(enemy), state)
			}
			onSelect={onSelect}
			value={value}
			title={t("simple.enemies")}
			placeholder={t("db.type_to_search")}
		/>
	);
}
