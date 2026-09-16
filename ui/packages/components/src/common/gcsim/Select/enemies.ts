import { dynamicKey, resources } from "@gcsim/localization";
import type { IEnemy } from "@gcsim/types";
import i18n from "i18next";
import { createKeyOrLabelPredicate } from "./utils";

export const enemies: IEnemy[] = Object.keys(resources.en.game.monster_names);

export function enemyLabel(enemy: IEnemy): string {
	return i18n.t(dynamicKey(`game:monster_names.${enemy}`));
}

export const enemyPredicate = createKeyOrLabelPredicate(enemyLabel);
