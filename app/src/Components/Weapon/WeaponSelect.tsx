import { Omnibar } from "@blueprintjs/select";
import { IWeapon, weaponSelectProps } from "./weapons";
import { useTranslation } from 'react-i18next'

const WeaponOmnibar = Omnibar.ofType<IWeapon>();

type Props = {
  isOpen: boolean;
  onClose: () => void;
  onSelect: (item: IWeapon) => void;
};

export function WeaponSelect(props: Props) {
  let { i18n } = useTranslation()

  return (
    <WeaponOmnibar
      resetOnSelect
      {...weaponSelectProps[i18n.language]}
      isOpen={props.isOpen}
      onClose={props.onClose}
      onItemSelect={props.onSelect}
    />
  );
}
