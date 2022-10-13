import { useTranslation } from "react-i18next";
import { SimResults } from "~src/Pages/Viewer/SimResults";
import SummaryCard from "../SummaryCard";

export default ({ data, color }: { data?: SimResults, color: string }) => {
  const { i18n } = useTranslation();
  const sps = data?.statistics?.sps;

  return (
    <SummaryCard
        key="sps"
        color={color}
        title="Shield HP Per Second (SPS)"
        value={sps?.mean?.toLocaleString(i18n.language, { maximumFractionDigits: 0 })}
        auxStats={[
          { title: "min", value: sps?.min?.toLocaleString(i18n.language, { maximumFractionDigits: 0 }) },
          { title: "max", value: sps?.max?.toLocaleString(i18n.language, { maximumFractionDigits: 0 }) },
          { title: "std", value: sps?.sd?.toLocaleString(i18n.language, { maximumFractionDigits: 0 }) },
          { title: "p25", value: sps?.q1?.toLocaleString(i18n.language, { maximumFractionDigits: 0 }) },
          { title: "p50", value: sps?.q2?.toLocaleString(i18n.language, { maximumFractionDigits: 0 }) },
          { title: "p75", value: sps?.q3?.toLocaleString(i18n.language, { maximumFractionDigits: 0 }) },
        ]}
        tooltip="help"
        drawerTitle="Shield Statistics">
      <div></div>
    </SummaryCard>
  );
};