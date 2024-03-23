import { CharacterBucketStats } from "@gcsim/types";
import { Group } from "@visx/group";
import { LegendItem, LegendLabel, LegendOrdinal } from "@visx/legend";
import { scaleLinear, scaleOrdinal } from "@visx/scale";
import { AreaStack, Bar, LinePath } from "@visx/shape";
import { DataColorsConst, useDataColors, GraphAxisBottom, GraphAxisLeft, GraphGrid, NoData } from "../../Util";
import { useTranslation } from "react-i18next";
import { useTooltip } from "@visx/tooltip";
import { HoverLine, RenderTooltip, TooltipData, useTooltipHandles } from "./CumulativeTooltip";
import { useData } from "./CumulativeData";
import { specialLocales } from "@ui/i18n";

const defaultMargin = { top: 10, left: 100, right: 20, bottom: 40 };

type LegendProps = {
  names?: string[];
}

export const CumulativeLegend = ({ names }: LegendProps) => {
  const { DataColors } = useDataColors();
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
                      fill={DataColorsConst.qualitative2(label.index)}
                      fillOpacity={0.5}
                      stroke={DataColorsConst.qualitative4(label.index)}
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

type GraphProps = {
  width: number;
  height: number;
  names?: string[];
  input?: CharacterBucketStats;
  margin?: { left: number, right: number, top: number, bottom: number };
}

export const CumulativeGraph = (
    {
      width,
      height,
      names,
      input,
      margin = defaultMargin
    }: GraphProps) => {
  const xMax = width - margin.left - margin.right;
  const yMax = height - margin.top - margin.bottom;
  const numXTicks = xMax < 475 ? 5 : 20;
  const numYTicks = 10;
  const bucketSize = input?.bucket_size ?? 30;

  const { i18n, t } = useTranslation();
  const { data, keys, duration } = useData(input, names);

  const xScale = scaleLinear<number>({
    range: [0, xMax],
    domain: [0, duration],
  });

  const yScale = scaleLinear<number>({
    range: [yMax, 0],
    domain: [0, 1],
  });

  const tooltip = useTooltip<TooltipData>();
  const tooltipHandles = useTooltipHandles(
      tooltip.showTooltip, tooltip.hideTooltip, xScale, yMax, margin, bucketSize);

  if (names == null || input == null || data.length == 0) {
    return <NoData />;
  }
  
  return (
    <div className="relative">
      <svg width={width} height={height}>
        <Group left={margin.left} top={margin.top}>
          <GraphGrid
              opacity={0.35}
              strokeDasharray="0 5 0"
              numTicksColumns={numXTicks}
              xScale={xScale}
              numTicksRows={numYTicks}
              yScale={yScale}
              width={xMax}
              height={yMax}
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
                      stroke={DataColorsConst.qualitative4(stack.key)}
                      strokeWidth={2}
                  />
                  <path
                      d={path(stack) || ''}
                      fill={DataColorsConst.qualitative2(stack.key)}
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
          <GraphAxisLeft
              hideTicks
              scale={yScale}
              axisLineClassName="stroke-2"
              tickFormat={s => s.toLocaleString(i18n.language, { style: "percent" })}
              numTicks={numYTicks}
              labelOffset={65}
              labelProps={specialLocales.includes(i18n.resolvedLanguage) ? { transform: `scale(1 1) translate(56 204)`, style: { writingMode: "vertical-lr" }, textAnchor: "middle" } : undefined}
              label={t<string>("result.cumu_contrib")}
          />
          <GraphAxisBottom
              hideTicks
              top={yMax}
              scale={xScale}
              axisLineClassName="stroke-2"
              tickFormat={s => s + t<string>("result.seconds_short")}
              numTicks={numXTicks}
              labelOffset={10}
              label={`${t<string>("result.dur_long")} (${t<string>("result.seconds")})`}
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