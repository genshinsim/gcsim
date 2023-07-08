import { FloatStat } from "@gcsim/types";
import { Group } from "@visx/group";
import { scaleBand, scaleLinear } from "@visx/scale";
import { BarStackHorizontal } from "@visx/shape";
import { StackKey } from "@visx/shape/lib/types";
import { useTooltip } from "@visx/tooltip";
import { useTranslation } from "react-i18next";
import { GraphAxisBottom, GraphAxisLeft } from "../Axes";
import { GraphGridColumns } from "../Grids";
import { HoverBoxPlot } from "./HoverBoxPlot";
import { RenderTooltip, TooltipData, useTooltipHandles } from "./Tooltip";

type Props<Datum, Key extends StackKey> = {
  width: number;
  height: number;
  margin?: { left: number, right: number, top: number, bottom: number };

  xDomain: number[];
  yDomain: string[];
  
  data: Datum[];
  keys: Key[];
  y: (d: Datum) => string;
  value: (d: Datum, k: Key) => number;
  stat: (d: Datum, k: Key) => FloatStat;

  barColor: (k: Key) => string;
  hoverColor: (k: Key) => string;

  tooltipContent?: (d: Datum, k: Key) => string | JSX.Element;
  bottomLabel?: string;
}

const defaultMargin = { top: 0, left: 80, right: 20, bottom: 40 };

export default <Datum,Key extends StackKey>(
    {
      width,
      height,
      margin = defaultMargin,
      data,
      keys,
      xDomain,
      yDomain,
      y,
      value,
      stat,
      barColor,
      hoverColor,
      tooltipContent,
      bottomLabel,
    }: Props<Datum,Key>) => {
  const tooltip = useTooltip<TooltipData<Key>>();
  const tooltipHandles = useTooltipHandles(tooltip.showTooltip, tooltip.hideTooltip);

  const { i18n } = useTranslation();
  const xMax = width - margin.left - margin.right;
  const yMax = height - margin.top - margin.bottom;

  const xScale = scaleLinear<number>({
    range: [0, xMax],
    domain: xDomain,
    nice: true
  });

  const yScale = scaleBand<string>({
    range: [yMax, 0],
    domain: yDomain,
    padding: 0.15,
  });

  return (
    <div className="relative">
      <svg width={width} height={height}>
        <Group left={margin.left} top={margin.top}>
          <GraphAxisLeft hideAxisLine hideTicks scale={yScale} numTicks={10000} />
          <GraphAxisBottom
              hideTicks
              hideAxisLine
              top={yMax}
              scale={xScale}
              label={bottomLabel}
              tickFormat={s => s.toLocaleString(
                  i18n.language, { notation: 'compact', maximumSignificantDigits: 3 })
              } />
          <GraphGridColumns scale={xScale} width={width} height={yMax} />
          <BarStackHorizontal<Datum, Key>
              data={data}
              keys={keys}
              height={yMax}
              y={y}
              xScale={xScale}
              yScale={yScale}
              color={barColor}
              value={value}>
            {(barStacks) =>
              barStacks.map((barStack =>
                barStack.bars.map((bar) => {
                  if (bar.width <= 0) {
                    return null;
                  }

                  const hoverData: TooltipData<Key> = {
                    index: bar.index,
                    key: bar.key,
                    x: bar.x,
                    y: bar.y,
                    height: bar.height,
                    width: bar.width,
                  };

                  const hover = (
                    bar.index === tooltip.tooltipData?.index
                    && bar.key === tooltip.tooltipData?.key
                  );
                  return (
                    <Group
                        key={"barstack-" + barStack.index + "-" + bar.index}
                        top={bar.y} left={bar.x}>
                      <rect
                          width={bar.width}
                          height={bar.height}
                          fill={hover ? hoverColor(bar.key) : bar.color}
                          stroke="#FFF"
                          strokeWidth={0.25}
                          strokeOpacity={1}
                          onMouseLeave={() => tooltipHandles.mouseLeave()}
                          onMouseMove={(e) => tooltipHandles.mouseHover(e, hoverData)} />
                    </Group>
                  );
                })))
            }
          </BarStackHorizontal>
          <HoverBoxPlot
              data={data}
              tooltip={tooltip.tooltipData}
              open={tooltip.tooltipOpen}
              scale={xScale}
              color={hoverColor}
              handles={tooltipHandles}
              stat={stat} />
        </Group>
      </svg>
      <RenderTooltip
          data={data}
          tooltipOpen={tooltip.tooltipOpen}
          tooltipData={tooltip.tooltipData}
          tooltipLeft={tooltip.tooltipLeft}
          tooltipTop={tooltip.tooltipTop}
          handles={tooltipHandles}
          showTooltip={tooltip.showTooltip}
          content={tooltipContent}
      />
    </div>
  );
};