import { SimResults } from "@gcsim/types";
import { useTranslation } from "react-i18next";
import SummaryCard from "../SummaryCard";

export default ({ data, color }: { data: SimResults | null; color: string }) => {
  const { i18n } = useTranslation();
  const fmt = (val?: number) => val?.toLocaleString(i18n.language, { maximumFractionDigits: 0 });
  const sps = data?.statistics?.sps;

  return (
    <SummaryCard
      key="sps"
      color={color}
      title="Shield HP Per Second (SPS)"
      value={fmt(sps?.mean)}
      auxStats={[
        { title: "min", value: fmt(sps?.min) },
        { title: "max", value: fmt(sps?.max) },
        { title: "std", value: fmt(sps?.sd) },
        { title: "p25", value: fmt(sps?.q1) },
        { title: "p50", value: fmt(sps?.q2) },
        { title: "p75", value: fmt(sps?.q3) },
      ]}
      tooltip="help"
      drawerTitle="Shield Statistics"
    >
      <div></div>
    </SummaryCard>
  );
};
