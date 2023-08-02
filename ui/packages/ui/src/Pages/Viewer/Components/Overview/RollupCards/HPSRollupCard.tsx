import { SimResults } from "@gcsim/types";
import { useCallback, useMemo } from "react";
import { useTranslation } from "react-i18next";
import { useRefresh } from "../../Util";
import { RollupCard } from "./Template";

export const HPSRollupCard = ({ data, color }: { data: SimResults | null; color: string }) => {
  const { i18n } = useTranslation();
  const fmt = useCallback(
    (val?: number) => val?.toLocaleString(i18n.language, { maximumFractionDigits: 0 }), [i18n]);

  const hps =  useRefresh(d => d?.statistics?.hps, 200, data);
  const auxStats = useMemo(() => [
    { title: "min", value: fmt(hps?.min) },
    { title: "max", value: fmt(hps?.max) },
    { title: "std", value: fmt(hps?.sd) },
    { title: "p25", value: fmt(hps?.q1) },
    { title: "p50", value: fmt(hps?.q2) },
    { title: "p75", value: fmt(hps?.q3) },
  ], [hps, fmt]);

  return (
    <RollupCard
        key="hps"
        color={color}
        title="Healing Per Second (HPS)"
        value={fmt(hps?.mean)}
        auxStats={auxStats}
        tooltip="help"
        hashLink="healing" />
  );
};
