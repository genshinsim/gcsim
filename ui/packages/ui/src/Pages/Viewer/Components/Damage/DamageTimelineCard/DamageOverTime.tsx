import { BucketStats } from "@gcsim/types";
import { curveBasis } from "@visx/curve";
import { Group } from "@visx/group";
import { LegendItem, LegendLabel, LegendOrdinal } from "@visx/legend";
import { scaleLinear, scaleOrdinal } from "@visx/scale";
import { Bar, LinePath } from "@visx/shape";
import { Threshold } from "@visx/threshold";
import { useTranslation } from "react-i18next";
import { DataColorsConst, GraphAxisBottom, GraphAxisLeft, GraphGrid, NoData } from "../../Util";
import { HoverLine, RenderTooltip, TooltipData, useTooltipHandles } from "./DamageOverTimeTooltip";
import { useData } from "./DamageOverTimeData";
import { useTooltip } from "@visx/tooltip";
import { useRef } from "react";
import { LegendGlyph } from ".";
import { specialLocales } from "@ui/Translation/i18n";

const defaultMargin = { top: 10, left: 100, right: 20, bottom: 40 };

type LegendProps = {
  names: string[];
  glyphs: LegendGlyph[];
}

export const DamageOverTimeLegend = ({names, glyphs}: LegendProps) => {
  const scale = scaleOrdinal({
    domain: names,
    range: glyphs,
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
                      fill={label.value?.fill}
                      fillOpacity={label.value?.fillOpacity}
                      stroke={label.value?.stroke}
                      strokeOpacity={label.value?.strokeOpacity}
                      strokeWidth={2}
                      width={glpyhSize}
                      height={glpyhSize}
                      strokeDasharray={label.value?.strokeDashArray}
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
  names: string[];
  input?: BucketStats;
  margin?: { left: number, right: number, top: number, bottom: number };
}

export const DamageOverTimeGraph = (
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

  const minRef = useRef<SVGPathElement | null>(null);
  const meanRef = useRef<SVGPathElement | null>(null);
  const maxRef = useRef<SVGPathElement | null>(null);

  const { i18n, t } = useTranslation();
  const { data, duration, maxValue } = useData(input);

  const xScale = scaleLinear<number>({
    range: [0, xMax],
    domain: [0, duration],
  });

  const yScale = scaleLinear<number>({
    range: [yMax, 0],
    domain: [0, maxValue],
    nice: true,
  });

  const tooltip = useTooltip<TooltipData>();
  const tooltipHandles = useTooltipHandles(
      tooltip.showTooltip, tooltip.hideTooltip, xScale, yMax, margin, bucketSize);

  if (input == null || data.length == 0) {
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
          <Threshold
            id={`${Math.random()}`}
            data={data}
            x={(d) => xScale(d.x)}
            y0={(d) => yScale((d.y.mean ?? 0) + (d.y.sd ?? 0))}
            y1={(d) => yScale(Math.max((d.y.mean ?? 0) - (d.y.sd ?? 0), 0))}
            clipAboveTo={0}
            clipBelowTo={yMax}
            curve={curveBasis}
            aboveAreaProps={{
              fill: DataColorsConst.qualitative1(0),
              fillOpacity: 0.2
            }}
          />
          <LinePath
            data={data}
            x={(d) => xScale(d.x)}
            y={(d) => yScale((d.y.mean ?? 0) + (d.y.sd ?? 0))}
            strokeWidth={1}
            stroke={DataColorsConst.qualitative3(0)}
            opacity={0.5}
            strokeDasharray="0 8 0"
            curve={curveBasis}
          />
          <LinePath
            data={data}
            x={(d) => xScale(d.x)}
            y={(d) => yScale(Math.max((d.y.mean ?? 0) - (d.y.sd ?? 0), 0))}
            strokeWidth={1}
            stroke={DataColorsConst.qualitative3(0)}
            opacity={0.5}
            strokeDasharray="0 8 0"
            curve={curveBasis}
          />
          <LinePath
            data={data}
            x={(d) => xScale(d.x)}
            y={(d) => yScale(d.y.max ?? 0)}
            strokeWidth={1}
            // stroke={DataColorsConst.gray}
            stroke={DataColorsConst.qualitative2(1)}
            curve={curveBasis}
            opacity={0.5}
            innerRef={maxRef}
          />
          <LinePath
            data={data}
            x={(d) => xScale(d.x)}
            y={(d) => yScale(d.y.min ?? 0)}
            strokeWidth={1}
            stroke={DataColorsConst.qualitative2(3)}
            curve={curveBasis}
            opacity={0.5}
            innerRef={minRef}
          />
          <LinePath
            data={data}
            x={(d) => xScale(d.x)}
            y={(d) => yScale(d.y.mean ?? 0)}
            strokeWidth={2}
            stroke={DataColorsConst.qualitative3(8)}
            curve={curveBasis}
            innerRef={meanRef}
          />
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
              tickFormat={s => s.toLocaleString(
                i18n.language, { notation: 'compact', maximumSignificantDigits: 3 })}
              numTicks={numYTicks}
              labelOffset={65}
              labelProps={specialLocales.includes(i18n.resolvedLanguage) ? { transform: "scale(1 1) translate(56 204)", style: { writingMode: "vertical-lr" }, textAnchor: "middle" } : undefined}
              label={t<string>("result.dmg_over_time")}
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
              minRef={minRef}
              meanRef={meanRef}
              maxRef={maxRef}
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