import { SimResults } from "@gcsim/types";
import { useTranslation } from "react-i18next";
import OverviewCard from "../OverviewCard";

export default ({ data, color }: { data: SimResults | null; color: string }) => {
  const { i18n } = useTranslation();
  const fmt = (val?: number) => val?.toLocaleString(i18n.language, { maximumFractionDigits: 2 });
  const eps = data?.statistics?.eps;

  return (
    <OverviewCard
      key="eps"
      color={color}
      title="Energy Per Second (EPS)"
      value={fmt(eps?.mean)}
      auxStats={[
        { title: "min", value: fmt(eps?.min) },
        { title: "max", value: fmt(eps?.max) },
        { title: "std", value: fmt(eps?.sd) },
        { title: "p25", value: fmt(eps?.q1) },
        { title: "p50", value: fmt(eps?.q2) },
        { title: "p75", value: fmt(eps?.q3) },
      ]}
      tooltip="help"
      drawerTitle="Energy Statistics"
    >
      <div></div>
    </OverviewCard>
  );
};
