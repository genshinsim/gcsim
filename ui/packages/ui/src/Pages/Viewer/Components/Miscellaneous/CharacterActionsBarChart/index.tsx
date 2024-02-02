import { Card } from "@blueprintjs/core";
import { SimResults } from "@gcsim/types";
import { ParentSize } from "@visx/responsive";
import { CardTitle, useDataColors, useRefreshWithTimer } from "../../Util";
import { BarChart, BarChartLegend } from "./BarChart";
import { useTranslation } from "react-i18next";

type Props = {
  data: SimResults | null;
  running: boolean;
  names?: string[];
};

export const CharacterActionsCard = ({ data, running, names }: Props) => {
  const { DataColors } = useDataColors();
  const { t } = useTranslation();
  const [stats, timer] = useRefreshWithTimer(
    (d) => {
      return {
        data: d?.statistics?.character_actions ? d?.statistics?.character_actions.map(
          (s) => s.sources ? { sources: Object.fromEntries(Object.entries(s.sources).map(([k, v]) => [t<string>("actions."+k), v])) } : {}
        ) : undefined,
      };
    },
    5000,
    data,
    running
  );

  // get actions used in config and sort them similarly to docs so that order is predictable across sims
  const actionNames = stats.data
    ? [
        ...new Set(
          stats.data.flatMap((n) => (n.sources ? Object.keys(n.sources) : []))
        ),
      ].sort((a, b) => DataColors.actionKeys.indexOf(a) - DataColors.actionKeys.indexOf(b))
    : null;

  return (
    <Card className="flex flex-col col-span-3 min-h-[384px]">
      <div className="flex flex-col sm:flex-row justify-start gap-5">
        <div className="flex flex-col gap-2">
          <CardTitle title={t<string>("simple.actions")} tooltip="x" timer={timer} />
        </div>
        <div className="flex flex-grow justify-start sm:justify-center pb-5 sm:pb-0 items-center">
          <BarChartLegend actionNames={actionNames} />
        </div>
      </div>
      <ParentSize className="flex-grow">
        {({ width, height }) => (
          <BarChart
            width={width}
            height={height}
            actions={stats.data}
            names={names}
            actionNames={actionNames}
          />
        )}
      </ParentSize>
    </Card>
  );
};
