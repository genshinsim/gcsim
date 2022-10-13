import { useTranslation } from "react-i18next";
import { SimResults } from "~src/Pages/Viewer/SimResults";
import SummaryCard from "../SummaryCard";

export default ({ data, color }: { data?: SimResults, color: string }) => {
  const { i18n } = useTranslation();
  const dps = data?.statistics?.dps;

  return (
    <SummaryCard
        key="dps"
        color={color}
        title="Damage Per Second (DPS)"
        value={dps?.mean?.toLocaleString(i18n.language, { maximumFractionDigits: 0 })}
        auxStats={[
          { title: "min", value: dps?.min?.toLocaleString(i18n.language, { maximumFractionDigits: 0 }) },
          { title: "max", value: dps?.max?.toLocaleString(i18n.language, { maximumFractionDigits: 0 }) },
          { title: "std", value: dps?.sd?.toLocaleString(i18n.language, { maximumFractionDigits: 0 }) },
          { title: "p25", value: dps?.q1?.toLocaleString(i18n.language, { maximumFractionDigits: 0 }) },
          { title: "p50", value: dps?.q2?.toLocaleString(i18n.language, { maximumFractionDigits: 0 }) },
          { title: "p75", value: dps?.q3?.toLocaleString(i18n.language, { maximumFractionDigits: 0 }) },
        ]}
        tooltip="help"
        drawerTitle="Damage Statistics">
      <div></div>
    </SummaryCard>
  );
};