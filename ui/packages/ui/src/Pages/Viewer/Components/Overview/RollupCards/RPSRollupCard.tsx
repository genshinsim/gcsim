import { SimResults } from "@gcsim/types";
import { useTranslation } from "react-i18next";
import { RollupCard } from "./Template";

export default ({ data, color }: { data: SimResults | null; color: string }) => {
  const { i18n } = useTranslation();
  const fmt = (val?: number) => val?.toLocaleString(i18n.language, { maximumFractionDigits: 2 });
  const rps = data?.statistics?.rps;

  return (
    <RollupCard
        key="rps"
        color={color}
        title="Reactions Per Second (RPS)"
        value={fmt(rps?.mean)}
        auxStats={[
          { title: "min", value: fmt(rps?.min) },
          { title: "max", value: fmt(rps?.max) },
          { title: "std", value: fmt(rps?.sd) },
          { title: "p25", value: fmt(rps?.q1) },
          { title: "p50", value: fmt(rps?.q2) },
          { title: "p75", value: fmt(rps?.q3) },
        ]}
        tooltip="help"
        drawerTitle="Reaction Statistics" />
  );
};
