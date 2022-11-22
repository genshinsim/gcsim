import { SimResults } from "@gcsim/types";
import { useTranslation } from "react-i18next";
import OverviewCard from "../OverviewCard";

export default ({ data, color }: { data: SimResults | null; color: string }) => {
  const { i18n } = useTranslation();
  const fmt = (val?: number) => val?.toLocaleString(i18n.language, { maximumFractionDigits: 0 });
  const hps = data?.statistics?.hps;

  return (
    <OverviewCard
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
      drawerTitle="Healing Statistics"
    >
      <div>
        sdasda sadasd asd waeawsdas dasd as dqaweasdasd awsedwqd asdas dawsd qwedsad ase wqdasd wa
        easdawe wadsawea sdsad wad sadawe awewasdsawe{" "}
      </div>
    </OverviewCard>
  );
};
