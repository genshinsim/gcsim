import { Tag } from "@blueprintjs/core";
import { EnergySettings } from "@gcsim/types";
import { memo } from "react";
import { useTranslation } from "react-i18next";

type Props = {
  energy?: EnergySettings;
};

export const Energy = memo(({ energy }: Props) => {
  const { i18n } = useTranslation();

  if (energy == null || energy.start == null || energy.end == null || !energy.active) {
    return null;
  }

  const startSec = (energy.start / 60).toLocaleString(i18n.language, { maximumFractionDigits: 1 });
  const endSec = (energy.end / 60).toLocaleString(i18n.language, { maximumFractionDigits: 1 });
  const amount = (energy.amount ?? 0).toLocaleString(i18n.language, { maximumFractionDigits: 1 });

  return (
    <Tag large minimal intent="primary">
      <div className="flex flex-row items-center gap-1 font-mono select-none">
        <div className="text-xs text-gray-400 pr-1">interval</div>
        <div className="text-sm">{amount + "p"}</div>
        <div className="text-xs text-gray-400">every</div>
        <div className="text-sm">{startSec + "s"}</div>
        <div className="text-xs text-gray-400">to</div>
        <div className="text-sm">{endSec + "s"}</div>
      </div>
    </Tag>
  );
});