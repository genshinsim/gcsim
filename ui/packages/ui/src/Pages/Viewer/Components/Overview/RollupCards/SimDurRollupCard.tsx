import { SimResults } from "@gcsim/types";
import { useCallback, useMemo } from "react";
import { useTranslation } from "react-i18next";
import { useRefresh } from "../../Util";
import { RollupCard } from "./Template";

export const SimDurRollupCard = ({ data, color }: { data: SimResults | null; color: string }) => {
  const { i18n, t } = useTranslation();
  const fmt = useCallback(
      (val?: number) => val?.toLocaleString(i18n.language, { maximumFractionDigits: 2 }), [i18n]);
  
  const duration =  useRefresh(d => d?.statistics?.duration, 200, data);
  const auxStats = useMemo(() => [
    { title: "min", value: fmt(duration?.min) },
    { title: "max", value: fmt(duration?.max) },
    { title: "std", value: fmt(duration?.sd) },
    { title: "p25", value: fmt(duration?.q1) },
    { title: "p50", value: fmt(duration?.q2) },
    { title: "p75", value: fmt(duration?.q3) },
  ], [duration, fmt]);

  return (
    <RollupCard
        key="duration"
        color={color}
        title={`${t("result.dur_long")} (Dur)`}
        label={t<string>("result.seconds_short")}
        value={fmt(duration?.mean)}
        auxStats={auxStats}
        tooltip="help"
        hashLink="sim" />
  );
};
