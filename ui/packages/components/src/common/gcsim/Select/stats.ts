import { dynamicKey } from "@gcsim/localization";
import type { IStat } from "@gcsim/types";
import i18n from "i18next";
import { createKeyOrLabelPredicate } from "./utils";

export const stats: IStat[] = [
	"hp",
	"hp%",
	"atk",
	"atk%",
	"def",
	"def%",
	"cr",
	"cd",
	"er",
	"heal",
	"em",
	"phys%",
	"pyro%",
	"electro%",
	"hydro%",
	"dendro%",
	"anemo%",
	"geo%",
	"cryo%",
];

export function statLabel(stat: IStat): string {
	return i18n.t(dynamicKey(`stats.${stat}`));
}

export const statPredicate = createKeyOrLabelPredicate(statLabel);
