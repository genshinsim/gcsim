import { SimResults } from "@gcsim/types";
import { useCallback, useMemo } from "react";
import { useTranslation } from "react-i18next";
import { useRefresh } from "../../Util";
import { RollupCard } from "./Template";

export const SHPRollupCard = ({ data, color }: { data: SimResults | null; color: string }) => {
  const { i18n } = useTranslation();
  const fmt = useCallback(
    (val?: number) => val?.toLocaleString(i18n.language, { maximumFractionDigits: 0 }), [i18n]);
  
  const shp =  useRefresh(d => d?.statistics?.shp, 200, data);
  const auxStats = useMemo(() => [
    { title: "min", value: fmt(shp?.min) },
    { title: "max", value: fmt(shp?.max) },
    { title: "std", value: fmt(shp?.sd) },
    { title: "p25", value: fmt(shp?.q1) },
    { title: "p50", value: fmt(shp?.q2) },
    { title: "p75", value: fmt(shp?.q3) },
  ], [shp, fmt]);

  return (
    <RollupCard
        key="shp"
        color={color}
        title="Effective Shield HP (SHP)"
        value={fmt(shp?.mean)}
        auxStats={auxStats}
        tooltip="help"
        hashLink="shields" />
  );
};