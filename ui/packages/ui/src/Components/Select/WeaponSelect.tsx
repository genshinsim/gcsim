import { Omnibar } from "@blueprintjs/select";
import { IWeapon } from "@gcsim/types";
import { useTranslation } from "react-i18next";
import { weaponSelectProps } from "./weapons";

const WeaponOmnibar = Omnibar.ofType<IWeapon>();

type Props = {
  isOpen: boolean;
  onClose: () => void;
  onSelect: (item: IWeapon) => void;
};

export function WeaponSelect(props: Props) {
  const { t } = useTranslation();
  return (
    <WeaponOmnibar
      resetOnSelect
      {...weaponSelectProps}
      isOpen={props.isOpen}
      onClose={props.onClose}
      onItemSelect={props.onSelect}
      inputProps={{ placeholder: `${t("db.type_to_search")}` }}
    />
  );
}
