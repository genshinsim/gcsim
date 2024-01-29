import { SimResults } from "@gcsim/types";
import { useCallback, useMemo } from "react";
import { useTranslation } from "react-i18next";
import { useRefresh } from "../../Util";
import { RollupCard } from "./Template";

export const EPSRollupCard = ({ data, color }: { data: SimResults | null; color: string }) => {
  const { i18n, t } = useTranslation();
  const fmt = useCallback(
    (val?: number) => val?.toLocaleString(i18n.language, { maximumFractionDigits: 2 }), [i18n]);

  const eps =  useRefresh(d => d?.statistics?.eps, 200, data);
  const auxStats = useMemo(() => [
    { title: "min", value: fmt(eps?.min) },
    { title: "max", value: fmt(eps?.max) },
    { title: "std", value: fmt(eps?.sd) },
    { title: "p25", value: fmt(eps?.q1) },
    { title: "p50", value: fmt(eps?.q2) },
    { title: "p75", value: fmt(eps?.q3) },
  ], [eps, fmt]);

  return (
    <RollupCard
        key="eps"
        color={color}
        title={`${t("result.eps_long")} (EPS)`}
        value={fmt(eps?.mean)}
        auxStats={auxStats}
        tooltip="help"
        hashLink="energy" />
  );
};
