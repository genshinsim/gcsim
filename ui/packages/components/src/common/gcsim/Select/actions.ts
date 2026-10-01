import { dynamicKey } from "@gcsim/localization";
import type { IAction } from "@gcsim/types";
import i18n from "i18next";

export const actions: IAction[] = [
	"attack",
	"charge",
	"aim",
	"skill",
	"burst",
	"low_plunge",
	"high_plunge",
	"dash",
	"jump",
	"walk",
	"swap",
];

export function actionLabel(action: IAction): string {
	return i18n.t(dynamicKey(`actions.${action}`));
}
