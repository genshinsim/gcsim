import { FloatStat } from "@gcsim/types";
import { AxisBottom, AxisLeft, TickRendererProps } from "@visx/axis";
import { Group } from "@visx/group";
import { LegendItem, LegendLabel, LegendOrdinal } from "@visx/legend";
import { scaleLinear, scaleOrdinal } from "@visx/scale";
import { AreaStack, Bar, LinePath } from "@visx/shape";
import { DataColors, NoData } from "../../Util";
import { GridColumns, GridRows } from "@visx/grid";
import { Text } from "@visx/text";
import { useTranslation } from "react-i18next";
import { useTooltip } from "@visx/tooltip";
import { HoverLine, RenderTooltip, TooltipData, useTooltipHandles } from "./CumulativeTooltip";
import { useData } from "./CumulativeData";

const defaultMargin = { top: 10, left: 80, right: 20, bottom: 40 };

type Props = {
  width: number;
  height: number;
  names?: string[];
  input?: FloatStat[][];
  bucketSize?: number;
  margin?: { left: number, right: number, top: number, bottom: number };
}

export const CumulativeLegend = ({ names }: { names?: string[] }) => {
  if (names == null) {
    return null;
  }

  const scale = scaleOrdinal({
    domain: names,
    range: names.map((n, i) => DataColors.character(i)),
  });

  const glpyhSize = 15;

  return (
    <LegendOrdinal scale={scale} direction="row" labelMargin="0 15px 0 0" className="flex-wrap">
      {(labels) => (
        <div className="flex flex-row">
          {labels.map((label, i) => (
            <LegendItem key={"legend-" + i}>
              <div className="my-[2px] mr-[4px]">
                <svg width={glpyhSize} height={glpyhSize}>
                  <rect
                      fill={DataColors.qualitative2(label.index)}
                      fillOpacity={0.5}
                      stroke={DataColors.qualitative4(label.index)}
                      strokeWidth={2}
                      width={glpyhSize}
                      height={glpyhSize}
                  />
                </svg>
              </div>
              <LegendLabel align="left" margin="0 15px 0 0">
                {label.text}
              </LegendLabel>
            </LegendItem>
          ))}
          
        </div>
      )}
    </LegendOrdinal>
  );
};

export const CumulativeGraph = (
    {
      width,
      height,
      names,
      input,
      bucketSize = 30,
      margin = defaultMargin
    }: Props) => {
  const xMax = width - margin.left - margin.right;
  const yMax = height - margin.top - margin.bottom;
  const numXTicks = xMax < 475 ? 5 : 20;
  const numYTicks = 10;

  const { i18n } = useTranslation();
  const { data, keys } = useData(input, bucketSize, names);

  const xScale = scaleLinear<number>({
    range: [0, xMax],
    domain: [0, Math.floor(((data.length-1) * bucketSize) / 60)],
  });

  const yScale = scaleLinear<number>({
    range: [yMax, 0],
    domain: [0, 1],
  });

  const tooltip = useTooltip<TooltipData>();
  const tooltipHandles = useTooltipHandles(
      tooltip.showTooltip, tooltip.hideTooltip, xScale, yMax, margin, bucketSize);

  if (names == null || input == null) {
    return <NoData />;
  }
  
  return (
    <div className="relative">
      <svg width={width} height={height}>
        <Group left={margin.left} top={margin.top}>
          <GridColumns
              scale={xScale}
              width={width}
              stroke={DataColors.gray}
              lineStyle={{ opacity: 0.35 }}
              strokeDasharray="0 5 0"
              numTicks={numXTicks}
              height={yMax}
          />
          <GridRows
              scale={yScale}
              height={height}
              width={xMax}
              stroke={DataColors.gray}
              strokeDasharray="0 5 0"
              lineStyle={{ opacity: 0.35 }}
              numTicks={numYTicks}
          />
          <AreaStack
              data={data}
              keys={keys}
              x={d => xScale(d.data.x)}
              value={(d, k) => (d.y[k].mean ?? 0)}
              y0={d => yScale(d[0])}
              y1={d => yScale(d[1])}
          >
            {({ stacks, path }) =>
              stacks.map((stack) => (
                <g key={"stack-" + stack.key}>
                  <LinePath
                      data={stack}
                      x={d => xScale(d.data.x)}
                      y={d => yScale(d[1])}
                      stroke={DataColors.qualitative4(stack.key)}
                      strokeWidth={2}
                  />
                  <path
                      d={path(stack) || ''}
                      fill={DataColors.qualitative2(stack.key)}
                      fillOpacity={0.5}
                  />
                </g>
              ))
            }
          </AreaStack>
          <Bar
              width={xMax}
              height={yMax}
              fill="transparent"
              onMouseMove={tooltipHandles.mouseHover}
              onMouseLeave={() => tooltipHandles.mouseLeave()}
          />
          <HoverLine
              data={data}
              xScale={xScale}
              yScale={yScale}
              yMax={yMax}
              tooltipData={tooltip.tooltipData}
              tooltipOpen={tooltip.tooltipOpen}
              tooltipLeft={tooltip.tooltipLeft}
              margin={margin}
          />
          <AxisLeft
              hideTicks
              scale={yScale}
              stroke={DataColors.gray}
              tickStroke={DataColors.gray}
              tickLineProps={{ opacity: 0.5 }}
              axisLineClassName="stroke-2"
              labelClassName="fill-gray-400 text-lg"
              tickClassName="fill-gray-400 font-mono text-xs"
              tickComponent={(props) => <TickLabel {...props} />}
              tickFormat={s =>
                s.toLocaleString(i18n.language, { style: "percent" })
              }
              numTicks={numYTicks}
              label="% Damage Contribution"
          />
          <AxisBottom
              hideTicks
              top={yMax}
              scale={xScale}
              numTicks={numXTicks}
              axisLineClassName="stroke-2"
              stroke={DataColors.gray}
              tickStroke={DataColors.gray}
              tickLineProps={{ opacity: 0.5 }}
              labelClassName="fill-gray-400 font-mono text-base"
              tickClassName="font-mono text-xs"
              tickLabelProps={() => ({
                fill: DataColors.gray,
                textAnchor: "middle"
              })}
              tickFormat={s => s + "s"}
              label="Duration (secs)"
          />
        </Group>
      </svg>
      <RenderTooltip
          names={names}
          data={data}
          tooltipOpen={tooltip.tooltipOpen}
          tooltipData={tooltip.tooltipData}
          tooltipLeft={tooltip.tooltipLeft}
          tooltipTop={tooltip.tooltipTop}
          handles={tooltipHandles}
          showTooltip={tooltip.showTooltip}
          margin={margin}
      />
    </div>
  );
};

const TickLabel = (props: TickRendererProps) => {
  return (
    <Text x={props.x} y={props.y} dx="-0.25em" dy="0.25em" textAnchor="end" className="cursor-default">
      {props.formattedValue}
    </Text>
  );
};