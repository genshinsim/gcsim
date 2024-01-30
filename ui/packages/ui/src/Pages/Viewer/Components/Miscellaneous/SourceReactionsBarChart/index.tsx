import { Card } from "@blueprintjs/core";
import { SimResults } from "@gcsim/types";
import { ParentSize } from "@visx/responsive";
import { CardTitle, useRefreshWithTimer } from "../../Util";
import { BarChart, BarChartLegend } from "./BarChart";
import { useTranslation } from "react-i18next";

type Props = {
  data: SimResults | null;
  running: boolean;
  names?: string[];
};

export const SourceReactionsCard = ({ data, running, names }: Props) => {
  const { t } = useTranslation();
  const [stats, timer] = useRefreshWithTimer(d => {
    return {
      data: d?.statistics?.source_reactions ? d?.statistics?.source_reactions.map(
        (s) => s.sources ? { sources: Object.fromEntries(Object.entries(s.sources).map(([k, v]) => [t<string>("reactions."+k), v])) } : {}
      ) : undefined,
    };
  }, 5000, data, running);

  return (
    <Card className="flex flex-col col-span-3 min-h-[384px]">
      <div className="flex flex-col sm:flex-row justify-start gap-5">
        <div className="flex flex-col gap-2">
          <CardTitle title={t<string>("result.per_source", { s: t<string>("result.reactions") })} tooltip="x" timer={timer} />
        </div>
        <div className="flex flex-grow justify-start sm:justify-center pb-5 sm:pb-0 items-center">
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