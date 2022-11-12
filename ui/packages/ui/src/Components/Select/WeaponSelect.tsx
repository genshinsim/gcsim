import { Omnibar } from "@blueprintjs/select";
import { IWeapon } from "@gcsim/types";
import { weaponSelectProps } from "./weapons";

const WeaponOmnibar = Omnibar.ofType<IWeapon>();

type Props = {
  isOpen: boolean;
  onClose: () => void;
  onSelect: (item: IWeapon) => void;
};

export function WeaponSelect(props: Props) {
  return (
    <WeaponOmnibar
      resetOnSelect
      {...weaponSelectProps}
      isOpen={props.isOpen}
      onClose={props.onClose}
      onItemSelect={props.onSelect}
      overlayProps={{ usePortal: false }}
    />
  );
}
