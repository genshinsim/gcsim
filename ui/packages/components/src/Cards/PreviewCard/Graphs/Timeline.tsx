import { model } from "@gcsim/types";
import { curveBasis } from "@visx/curve";
import { Group } from "@visx/group";
import { ParentSize } from "@visx/responsive";
import { scaleLinear } from "@visx/scale";
import { AreaClosed } from "@visx/shape";
import { DataColorsConst, NoDataIcon } from "../../../common/gcsim";
import { useData } from "../../ResultCards/Damage/DamageTimelineCard/DamageOverTimeData";

type Props = {
  data?: model.BucketStats | null;
};

export const Timeline = ({ data }: Props) => {
  if (data == null) {
    return <NoDataIcon className="h-16" />;
  }

  return (
    <ParentSize>
      {({ width, height }) => (
        <Graph width={width} height={height} input={data} />
      )}
    </ParentSize>
  );
};

type GraphProps = {
  width: number;
  height: number;
  input: model.BucketStats;
  margin?: { left: number; right: number; top: number; bottom: number };
};

const defaultMargin = { left: 8, right: 8, top: 0, bottom: 1 };

const Graph = ({
  width,
  height,
  input,
  margin = defaultMargin,
}: GraphProps) => {
  const xMax = width - margin.left - margin.right;
  const yMax = height - margin.top - margin.bottom;
  const { data, duration } = useData(input);
  const maxValue = data.reduce((acc, p) => Math.max(acc, p.y.mean ?? 0), 0.0);

  const xScale = scaleLinear<number>({
    range: [0, xMax],
    domain: [0, duration],
  });

  const yScale = scaleLinear<number>({
    range: [yMax, 0],
    domain: [0, maxValue],
  });

  const color = DataColorsConst.qualitative3(8);

  return (
    <div className="relative">
      <svg width={width} height={height}>
        <Group left={margin.left} top={margin.top}>
          <AreaClosed
            data={data}
            x={(d) => xScale(d.x)}
            y={(d) => yScale(d.y.mean ?? 0)}
            yScale={yScale}
            strokeWidth={2}
            stroke={color}
            fill={color}
            curve={curveBasis}
          />
        </Group>
      </svg>
    </div>
  );
};
