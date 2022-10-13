import { useTranslation } from "react-i18next";
import { SimResults } from "~src/Pages/Viewer/SimResults";
import SummaryCard from "../SummaryCard";


export default ({ data, color }: { data?: SimResults, color: string }) => {
  const { i18n } = useTranslation();
  const duration = data?.statistics?.duration;

  return (
    <SummaryCard
        key="duration"
        color={color}
        title="Sim Duration"
        label="s"
        value={duration?.mean?.toLocaleString(i18n.language, { maximumFractionDigits: 2 })}
        auxStats={[
          { title: "min", value: duration?.min?.toLocaleString(i18n.language, { maximumFractionDigits: 2 }) },
          { title: "max", value: duration?.max?.toLocaleString(i18n.language, { maximumFractionDigits: 2 }) },
          { title: "std", value: duration?.sd?.toLocaleString(i18n.language, { maximumFractionDigits: 2 }) },
          { title: "p25", value: duration?.q1?.toLocaleString(i18n.language, { maximumFractionDigits: 2 }) },
          { title: "p50", value: duration?.q2?.toLocaleString(i18n.language, { maximumFractionDigits: 2 }) },
          { title: "p75", value: duration?.q3?.toLocaleString(i18n.language, { maximumFractionDigits: 2 }) },
        ]}
        tooltip="help"
        drawerTitle="Sim Duration Statistics">
      <div></div>
    </SummaryCard>
  );
};