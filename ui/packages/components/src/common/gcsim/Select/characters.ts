import { dynamicKey, resources } from "@gcsim/localization";
import type { ICharacter } from "@gcsim/types";
import i18n from "i18next";

export const characters: ICharacter[] = Object.keys(
	resources.en.game.character_names,
);

export function characterLabel(character: ICharacter): string {
	return i18n.t(dynamicKey(`game:character_names.${character}`));
}
