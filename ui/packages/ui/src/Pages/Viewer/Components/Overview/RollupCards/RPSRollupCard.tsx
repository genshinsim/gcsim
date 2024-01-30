import { SimResults } from "@gcsim/types";
import { useCallback, useMemo } from "react";
import { useTranslation } from "react-i18next";
import { useRefresh } from "../../Util";
import { RollupCard } from "./Template";

export const RPSRollupCard = ({ data, color }: { data: SimResults | null; color: string }) => {
  const { i18n, t } = useTranslation();
  const fmt = useCallback(
    (val?: number) => val?.toLocaleString(i18n.language, { maximumFractionDigits: 2 }), [i18n]);

  const rps =  useRefresh(d => d?.statistics?.rps, 200, data);
  const auxStats = useMemo(() => [
    { title: "min", value: fmt(rps?.min) },
    { title: "max", value: fmt(rps?.max) },
    { title: "std", value: fmt(rps?.sd) },
    { title: "p25", value: fmt(rps?.q1) },
    { title: "p50", value: fmt(rps?.q2) },
    { title: "p75", value: fmt(rps?.q3) },
  ], [rps, fmt]);

  return (
    <RollupCard
        key="rps"
        color={color}
        title={`${t("result.rps_long")} (RPS)`}
        value={fmt(rps?.mean)}
        auxStats={auxStats}
        tooltip="help"
        hashLink="reactions" />
  );
};
