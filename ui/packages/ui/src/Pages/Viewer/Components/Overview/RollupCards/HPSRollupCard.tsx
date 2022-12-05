import { SimResults } from "@gcsim/types";
import { useTranslation } from "react-i18next";
import { RollupCard } from "./Template";

export default ({ data, color }: { data: SimResults | null; color: string }) => {
  const { i18n } = useTranslation();
  const fmt = (val?: number) => val?.toLocaleString(i18n.language, { maximumFractionDigits: 0 });
  const hps = data?.statistics?.hps;

  return (
    <RollupCard
        key="hps"
        color={color}
        title="Healing Per Second (HPS)"
        value={fmt(hps?.mean)}
        auxStats={[
          { title: "min", value: fmt(hps?.min) },
          { title: "max", value: fmt(hps?.max) },
          { title: "std", value: fmt(hps?.sd) },
          { title: "p25", value: fmt(hps?.q1) },
          { title: "p50", value: fmt(hps?.q2) },
          { title: "p75", value: fmt(hps?.q3) },
        ]}
        tooltip="help"
        hashLink="healing" />
  );
};
