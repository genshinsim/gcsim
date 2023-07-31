import { Card } from "@blueprintjs/core";
import { SimResults } from "@gcsim/types";
import { ParentSize } from "@visx/responsive";
import { CardTitle, useRefreshWithTimer } from "../../Util";
import { BarChart, BarChartLegend } from "./BarChart";

type Props = {
  data: SimResults | null;
  running: boolean;
  names?: string[];
};

export const SourceReactionsCard = ({ data, running, names }: Props) => {
  const [stats, timer] = useRefreshWithTimer(d => {
    return {
      data: d?.statistics?.source_reactions,
    };
  }, 5000, data, running);

  return (
    <Card className="flex flex-col col-span-3 h-96">
      <div className="flex flex-row justify-start gap-5">
        <div className="flex flex-col gap-2">
          <CardTitle title="Reactions Per Source" tooltip="x" timer={timer} />
        </div>
        <div className="flex flex-grow justify-center items-center">
            <BarChartLegend names={names} />
        </div>
      </div>
      <ParentSize>
        {({ width, height }) => (
          <BarChart width={width} height={height} reactions={stats.data} names={names} />
        )}
      </ParentSize>
    </Card>
  );
};