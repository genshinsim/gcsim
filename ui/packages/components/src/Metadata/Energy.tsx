import { Badge } from "@gcsim/primitives";
import type { EnergySettings } from "@gcsim/types";
import { memo } from "react";
import { Trans, useTranslation } from "react-i18next";

type Props = {
	energy?: EnergySettings;
};

export const Energy = memo(({ energy }: Props) => {
	const { i18n, t } = useTranslation();

	if (
		energy == null ||
		energy.start == null ||
		energy.end == null ||
		!energy.active
	) {
		return null;
	}

	const startSec = (energy.start / 60).toLocaleString(i18n.language, {
		maximumFractionDigits: 1,
	});
	const endSec = (energy.end / 60).toLocaleString(i18n.language, {
		maximumFractionDigits: 1,
	});
	const amount = (energy.amount ?? 0).toLocaleString(i18n.language, {
		maximumFractionDigits: 1,
	});

	return (
		<Badge variant="default">
			<span className="flex flex-row items-center gap-1 font-mono select-none">
				<Trans i18nKey="result.metadata_energy">
					<span className="text-xs text-gray-400 pr-1" />
					<span className="text-sm">{{ p: amount + "p" } as never}</span>
					<span className="text-xs text-gray-400" />
					<span className="text-sm">
						{{ s: startSec + t("result.seconds_short") } as never}
					</span>
					<span className="text-xs text-gray-400" />
					<span className="text-sm">
						{{ e: endSec + t("result.seconds_short") } as never}
					</span>
				</Trans>
			</span>
		</Badge>
	);
});
