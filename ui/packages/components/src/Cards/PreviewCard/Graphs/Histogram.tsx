import { model } from "@gcsim/types";
import { NoDataIcon } from "@gcsim/ui/src/Pages/Viewer/Components/Util/NoData";
import { Group } from "@visx/group";
import { ParentSize } from "@visx/responsive";
import { Colors } from "../../../common/gcsim";
import { useScales } from "../../DistributionCard/HistogramGraph";

type Props = {
  width: number;
  height: number;
  margin?: { left: number; right: number; top: number; bottom: number };

  data: model.OverviewStats;
  hist: number[];
  barColor?: string;
  accentColor?: string;
};

const defaultMargin = { left: 8, right: 8, top: 8, bottom: 1 };

export const Histogram = ({ data }: { data?: model.OverviewStats }) => {
  if (data?.histogram == null) {
    return <NoDataIcon className="h-16" />;
  }

  return (
    <ParentSize>
      {({ width, height }) => (
        <Graph
          width={width}
          height={height}
          data={data}
          // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
          hist={data.histogram!}
        />
      )}
    </ParentSize>
  );
};

export const Graph = ({
  width,
  height,
  margin = defaultMargin,
  data,
  hist,
  barColor = Colors.VERMILION3,
  accentColor = Colors.VERMILION1,
}: Props) => {
  const xMax = width - margin.left - margin.right;
  const yMax = height - margin.top - margin.bottom;
  const { xScale, yScale, delta } = useScales(data, xMax, yMax);

  return (
    <div className="relative">
      <svg width={width} height={height}>
        <Group left={margin.left} top={margin.top}>
          {hist.map((c, i) => {
            const barWidth = xScale.bandwidth();
            const barHeight = yMax - yScale(c);
            const barX = xScale(i) ?? 0;
            const barY = yMax - barHeight;

            if (
              c <= 0 ||
              barHeight < 0 ||
              data.mean == null ||
              data.min == null
            ) {
              return null;
            }

            let fill = barColor;
            if (i === Math.floor((delta ?? 0) * (data.mean - data.min))) {
              fill = accentColor;
            }

            return (
              <rect
                key={"bin-" + i}
                fill={fill}
                x={barX}
                y={barY}
                width={barWidth}
                height={barHeight}
              />
            );
          })}
        </Group>
      </svg>
    </div>
  );
};
