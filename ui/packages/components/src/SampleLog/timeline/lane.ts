import { dynamicKey } from "@gcsim/localization";
import { useCallback } from "react";
import { useTranslation } from "react-i18next";
import type { Lane } from "./model";

export function laneColor(lane: Lane): string {
	if (lane.element == null) {
		return "var(--g-text-mute)";
	}
	return `var(--g-${lane.element})`;
}

export function useLaneName() {
	const { t } = useTranslation();
	return useCallback(
		(lane: Lane) =>
			lane.index === 0
				? t("sample.sim_lane")
				: t(dynamicKey(`game:character_names.${lane.key}`)),
		[t],
	);
}
