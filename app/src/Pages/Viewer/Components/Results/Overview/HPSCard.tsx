import { useTranslation } from "react-i18next";
import { SimResults } from "~src/Pages/Viewer/SimResults";
import SummaryCard from "../SummaryCard";

export default ({ data, color }: { data?: SimResults, color: string }) => {
  const { i18n } = useTranslation();
  const hps = data?.statistics?.hps;

  return (
    <SummaryCard
        key="hps"
        color={color}
        title="Healing Per Second (HPS)"
        value={hps?.mean?.toLocaleString(i18n.language, { maximumFractionDigits: 2 })}
        auxStats={[
          { title: "min", value: hps?.min?.toLocaleString(i18n.language, { maximumFractionDigits: 2 }) },
          { title: "max", value: hps?.max?.toLocaleString(i18n.language, { maximumFractionDigits: 2 }) },
          { title: "std", value: hps?.sd?.toLocaleString(i18n.language, { maximumFractionDigits: 2 }) },
          { title: "p25", value: hps?.q1?.toLocaleString(i18n.language, { maximumFractionDigits: 2 }) },
          { title: "p50", value: hps?.q2?.toLocaleString(i18n.language, { maximumFractionDigits: 2 }) },
          { title: "p75", value: hps?.q3?.toLocaleString(i18n.language, { maximumFractionDigits: 2 }) },
        ]}
        tooltip="help"
        drawerTitle="Healing Statistics">
      <div>sdasda sadasd asd waeawsdas dasd as dqaweasdasd awsedwqd asdas dawsd qwedsad ase wqdasd wa easdawe wadsawea sdsad wad sadawe awewasdsawe </div>
    </SummaryCard>
  );
};