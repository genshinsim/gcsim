import { useTranslation } from "react-i18next";
import { SimResults } from "~src/Pages/Viewer/SimResults";
import SummaryCard from "../SummaryCard";

export default ({ data, color }: { data?: SimResults, color: string }) => {
  const { i18n } = useTranslation();
  const eps = data?.statistics?.eps;

  return (
    <SummaryCard
        key="eps"
        color={color}
        title="Energy Per Second (EPS)"
        value={eps?.mean?.toLocaleString(i18n.language, { maximumFractionDigits: 2 })}
        auxStats={[
          { title: "min", value: eps?.min?.toLocaleString(i18n.language, { maximumFractionDigits: 2 }) },
          { title: "max", value: eps?.max?.toLocaleString(i18n.language, { maximumFractionDigits: 2 }) },
          { title: "std", value: eps?.sd?.toLocaleString(i18n.language, { maximumFractionDigits: 2 }) },
          { title: "p25", value: eps?.q1?.toLocaleString(i18n.language, { maximumFractionDigits: 2 }) },
          { title: "p50", value: eps?.q2?.toLocaleString(i18n.language, { maximumFractionDigits: 2 }) },
          { title: "p75", value: eps?.q3?.toLocaleString(i18n.language, { maximumFractionDigits: 2 }) },
        ]}
        tooltip="help"
        drawerTitle="Energy Statistics">
      <div></div>
    </SummaryCard>
  );
};