import { Card } from "@blueprintjs/core";
import { SimResults } from "@gcsim/types";
import { ParentSize } from "@visx/responsive";
import { CardTitle, DataColors, useRefreshWithTimer } from "../../Util";
import { BarChart, BarChartLegend } from "./BarChart";

type Props = {
  data: SimResults | null;
  running: boolean;
  names?: string[];
};

export const CharacterActionsCard = ({ data, running, names }: Props) => {
  const [stats, timer] = useRefreshWithTimer(
    (d) => {
      return {
        data: d?.statistics?.character_actions,
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
    <Card className="flex flex-col col-span-full h-96">
      <div className="flex flex-row justify-start gap-5">
        <div className="flex flex-col gap-2">
          <CardTitle title="Character Actions" tooltip="x" timer={timer} />
        </div>
        <div className="flex flex-grow justify-center items-center">
          <BarChartLegend actionNames={actionNames} />
        </div>
      </div>
      <ParentSize>
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
