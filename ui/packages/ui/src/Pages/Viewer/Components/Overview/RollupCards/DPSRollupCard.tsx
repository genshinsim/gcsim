import { SimResults } from "@gcsim/types";
import { useCallback, useMemo } from "react";
import { useTranslation } from "react-i18next";
import { useRefresh } from "../../Util";
import { RollupCard } from "./Template";

export const DPSRollupCard = ({ data, color }: { data: SimResults | null; color: string }) => {
  const { i18n, t } = useTranslation();
  const fmt = useCallback(
    (val?: number) => val?.toLocaleString(i18n.language, { maximumFractionDigits: 0 }), [i18n]);
  
  const dps = useRefresh(d => d?.statistics?.dps, 200, data);
  const auxStats = useMemo(() => [
    { title: "min", value: fmt(dps?.min) },
    { title: "max", value: fmt(dps?.max) },
    { title: "std", value: fmt(dps?.sd) },
    { title: "p25", value: fmt(dps?.q1) },
    { title: "p50", value: fmt(dps?.q2) },
    { title: "p75", value: fmt(dps?.q3) },
  ], [dps, fmt]);

  return (
    <RollupCard
        key="dps"
        color={color}
        title={`${t("result.dps_long")} (DPS)`}
        value={fmt(dps?.mean)}
        auxStats={auxStats}
        tooltip="help"
        hashLink="damage" />
  );
};
