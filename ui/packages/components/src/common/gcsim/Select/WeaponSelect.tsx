import type { IWeapon } from "@gcsim/types";
import { useTranslation } from "react-i18next";
import { OmniSelect } from "./OmniSelect";
import { renderLabeledItem } from "./utils";
import { weaponLabel, weaponPredicate, weapons } from "./weapons";

export type WeaponSelectProps = {
	isOpen: boolean;
	onClose: () => void;
	onSelect: (weapon: IWeapon) => void;
	value?: IWeapon;
};

export function WeaponSelect({
	isOpen,
	onClose,
	onSelect,
	value,
}: WeaponSelectProps) {
	const { t } = useTranslation();
	return (
		<OmniSelect<IWeapon>
			isOpen={isOpen}
			onClose={onClose}
			items={weapons}
			itemKey={(weapon) => weapon}
			itemPredicate={weaponPredicate}
			itemRenderer={(weapon, state) =>
				renderLabeledItem(weapon, weaponLabel(weapon), state)
			}
			onSelect={onSelect}
			value={value}
			title={t("simple.weapons")}
			placeholder={t("db.type_to_search")}
		/>
	);
}
