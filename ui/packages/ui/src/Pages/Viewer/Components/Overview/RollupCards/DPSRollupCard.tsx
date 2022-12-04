import { SimResults } from "@gcsim/types";
import { useTranslation } from "react-i18next";
import { RollupCard } from "./Template";

export default ({ data, color }: { data: SimResults | null; color: string }) => {
  const { i18n } = useTranslation();
  const fmt = (val?: number) => val?.toLocaleString(i18n.language, { maximumFractionDigits: 0 });
  const dps = data?.statistics?.dps;

  return (
    <RollupCard
        key="dps"
        color={color}
        title="Damage Per Second (DPS)"
        value={fmt(dps?.mean)}
        auxStats={[
          { title: "min", value: fmt(dps?.min) },
          { title: "max", value: fmt(dps?.max) },
          { title: "std", value: fmt(dps?.sd) },
          { title: "p25", value: fmt(dps?.q1) },
          { title: "p50", value: fmt(dps?.q2) },
          { title: "p75", value: fmt(dps?.q3) },
        ]}
        tooltip="help"
        drawerTitle="Damage Statistics" />
  );
};
