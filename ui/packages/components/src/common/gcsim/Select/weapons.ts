import { dynamicKey, resources } from "@gcsim/localization";
import type { IWeapon } from "@gcsim/types";
import i18n from "i18next";
import { createKeyOrLabelPredicate } from "./utils";

export const weapons: IWeapon[] = Object.keys(resources.en.game.weapon_names);

export function weaponLabel(weapon: IWeapon): string {
	return i18n.t(dynamicKey(`game:weapon_names.${weapon}`));
}

export const weaponPredicate = createKeyOrLabelPredicate(weaponLabel);
