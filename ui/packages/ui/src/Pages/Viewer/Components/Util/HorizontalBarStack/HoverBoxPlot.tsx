import { TooltipData, TooltipHandles } from "./Tooltip";
import { ScaleLinear } from "d3-scale";
import { Group } from "@visx/group";
import { BoxPlot } from "@visx/stats";
import { FloatStat } from "@gcsim/types";

type Props<Datum,Key> = {
  data: Datum[];
  tooltip?: TooltipData<Key>;
  open: boolean;
  scale: ScaleLinear<number, number>;
  color: (k: Key) => string;
  handles: TooltipHandles<Key>;
  stat: (d: Datum, k: Key) => FloatStat;
}

export const HoverBoxPlot = <Datum,Key>(
    { data, tooltip, scale, open, color, handles, stat }: Props<Datum,Key>) => {
  if (!tooltip || !open) {
    return null;
  }

  const value = stat(data[tooltip.index], tooltip.key);

  return (
    <Group
        left={tooltip.x}
        top={tooltip.y}
        onMouseLeave={() => handles.mouseLeave()}
        onMouseMove={(e) => handles.mouseHover(e, tooltip)}
    >
      <BoxPlot
          horizontal
          top={10}
          boxWidth={tooltip.height - 20}
          valueScale={scale}
          min={value.min}
          max={value.max}
          firstQuartile={(value.mean ?? 0) - (value.sd ?? 0)}
          median={value.mean}
          thirdQuartile={(value.mean ?? 0) + (value.sd ?? 0)}
          stroke={"#FFF"}
          fill={color(tooltip.key)}
      />
    </Group>
  );
};