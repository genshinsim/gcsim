import { useTranslation } from "react-i18next";
import { SimResults } from "~src/Pages/Viewer/SimResults";
import SummaryCard from "../SummaryCard";

export default ({ data, color }: { data?: SimResults, color: string }) => {
  const { i18n } = useTranslation();
  const rps = data?.statistics?.rps;

  return (
    <SummaryCard
        key="rps"
        color={color}
        title="Reactions Per Second (RPS)"
        value={rps?.mean?.toLocaleString(i18n.language, { maximumFractionDigits: 2 })}
        auxStats={[
          { title: "min", value: rps?.min?.toLocaleString(i18n.language, { maximumFractionDigits: 2 }) },
          { title: "max", value: rps?.max?.toLocaleString(i18n.language, { maximumFractionDigits: 2 }) },
          { title: "std", value: rps?.sd?.toLocaleString(i18n.language, { maximumFractionDigits: 2 }) },
          { title: "p25", value: rps?.q1?.toLocaleString(i18n.language, { maximumFractionDigits: 2 }) },
          { title: "p50", value: rps?.q2?.toLocaleString(i18n.language, { maximumFractionDigits: 2 }) },
          { title: "p75", value: rps?.q3?.toLocaleString(i18n.language, { maximumFractionDigits: 2 }) },
        ]}
        tooltip="help"
        drawerTitle="Reaction Statistics">
      <div></div>
    </SummaryCard>
  );
};