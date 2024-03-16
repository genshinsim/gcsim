import { Card } from "@blueprintjs/core";
import { SimResults } from "@gcsim/types";
import { ParentSize } from "@visx/responsive";
import { CardTitle, useRefreshWithTimer } from "../../Util";
import { BarChart } from "./BarChart";
import { useTranslation } from "react-i18next";

type Props = {
  data: SimResults | null;
  running: boolean;
  names?: string[];
};

export const EndingEnergyCard = ({ data, running, names }: Props) => {
  const { t } = useTranslation();
  const [stats, timer] = useRefreshWithTimer(d => {
    return {
      data: d?.statistics?.end_stats,
    };
  }, 5000, data, running);

  return (
    <Card className="flex flex-col col-span-full h-auto">
      <div className="flex flex-col sm:flex-row justify-start gap-5">
        <div className="flex flex-col gap-2">
          <CardTitle title={t<string>("result.ending_energy")} tooltip="x" timer={timer} />
        </div>
      </div>
      <ParentSize>
        {({ width, height }) => (
          <BarChart width={width} height={height} end_stats={stats.data} names={names} />
        )}
      </ParentSize>
    </Card>
  );
};