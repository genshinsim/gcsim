import { BucketStats } from "@gcsim/types";
import { Group } from "@visx/group";
import { LegendItem, LegendLabel, LegendOrdinal } from "@visx/legend";
import { scaleLinear, scaleOrdinal } from "@visx/scale";
import { Bar, LinePath } from "@visx/shape";
import { useTooltip } from "@visx/tooltip";
import { useRef } from "react";
import { useTranslation } from "react-i18next";
import {
  DataColors,
  GraphAxisBottom,
  GraphAxisLeft,
  GraphGrid,
  NoData,
} from "../../Util";
import { useData } from "./CumulativeDamageData";
import {
  HoverLine,
  RenderTooltip,
  TooltipData,
  useTooltipHandles,
} from "./CumulativeDamageTooltip";

const defaultMargin = { top: 10, left: 100, right: 20, bottom: 40 };

type Props = {
  width: number;
  height: number;
  graph: string;
  target: string;
  input?: BucketStats;
  margin?: { left: number; right: number; top: number; bottom: number };
};

type LegendGlyph = {
  label: string;
  fill: string;
  fillOpacity: number;
  stroke: string;
  strokeOpacity: number;
  strokeDashArray?: string;
};

const names = ["min", "max", "p25", "p50", "p75"];
const glyphs: LegendGlyph[] = [
  {
    label: "min",
    fill: DataColors.qualitative2(3),
    fillOpacity: 0.75,
    stroke: DataColors.qualitative2(3),
    strokeOpacity: 0.75,
    strokeDashArray: "0 5 0",
  },
  {
    label: "max",
    fill: DataColors.qualitative2(1),
    fillOpacity: 0.75,
    stroke: DataColors.qualitative2(1),
    strokeOpacity: 0.75,
    strokeDashArray: "0 5 0",
  },
  {
    label: "p25",
    fill: DataColors.qualitative2(4),
    fillOpacity: 1.0,
    stroke: DataColors.qualitative2(4),
    strokeOpacity: 0,
  },
  {
    label: "p50",
    fill: DataColors.qualitative3(8),
    fillOpacity: 1.0,
    stroke: DataColors.qualitative3(8),
    strokeOpacity: 0,
  },
  {
    label: "p75",
    fill: DataColors.qualitative2(5),
    fillOpacity: 1.0,
    stroke: DataColors.qualitative2(5),
    strokeOpacity: 0,
  },
];

export const CumulativeLegend = () => {
  const scale = scaleOrdinal({
    domain: names,
    range: glyphs,
  });

  const glpyhSize = 15;

  return (
    <LegendOrdinal
      scale={scale}
      direction="row"
      labelMargin="0 15px 0 0"
      className="flex-wrap"
    >
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

export const CumulativeGraph = ({
  width,
  height,
  graph,
  target,
  input,
  margin = defaultMargin,
}: Props) => {
  const xMax = width - margin.left - margin.right;
  const yMax = height - margin.top - margin.bottom;
  const numXTicks = xMax < 475 ? 5 : 20;
  const numYTicks = 10;
  const bucketSize = input?.bucket_size ?? 30;

  const minRef = useRef<SVGPathElement | null>(null);
  const maxRef = useRef<SVGPathElement | null>(null);
  const q1Ref = useRef<SVGPathElement | null>(null);
  const q2Ref = useRef<SVGPathElement | null>(null);
  const q3Ref = useRef<SVGPathElement | null>(null);

  const { i18n } = useTranslation();
  const { data, duration, maxValue } = useData(graph, target, input);

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
    tooltip.showTooltip,
    tooltip.hideTooltip,
    xScale,
    yMax,
    margin,
    bucketSize
  );

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
          <LinePath
            data={data}
            x={(d) => xScale(d.x)}
            y={(d) => yScale(d.y.min ?? 0)}
            strokeWidth={1}
            stroke={DataColors.qualitative2(3)}
            strokeDasharray="0 8 0"
            opacity={0.75}
            innerRef={minRef}
          />
          <LinePath
            data={data}
            x={(d) => xScale(d.x)}
            y={(d) => yScale(d.y.max ?? 0)}
            strokeWidth={1}
            stroke={DataColors.qualitative2(1)}
            strokeDasharray="0 8 0"
            opacity={0.75}
            innerRef={maxRef}
          />
          <LinePath
            data={data}
            x={(d) => xScale(d.x)}
            y={(d) => yScale(d.y.q1 ?? 0)}
            strokeWidth={2}
            stroke={DataColors.qualitative2(4)}
            innerRef={q1Ref}
          />
          <LinePath
            data={data}
            x={(d) => xScale(d.x)}
            y={(d) => yScale(d.y.q2 ?? 0)}
            strokeWidth={2}
            stroke={DataColors.qualitative3(8)}
            innerRef={q2Ref}
          />
          <LinePath
            data={data}
            x={(d) => xScale(d.x)}
            y={(d) => yScale(d.y.q3 ?? 0)}
            strokeWidth={2}
            stroke={DataColors.qualitative2(5)}
            innerRef={q3Ref}
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
            tickFormat={(s) =>
              s.toLocaleString(i18n.language, {
                notation: "compact",
                maximumSignificantDigits: 3,
              })
            }
            numTicks={numYTicks}
            labelOffset={65}
            label="Cumulative Damage"
          />
          <GraphAxisBottom
            hideTicks
            top={yMax}
            scale={xScale}
            axisLineClassName="stroke-2"
            tickFormat={(s) => s + "s"}
            numTicks={numXTicks}
            label="Duration (secs)"
          />
          <HoverLine
            data={data}
            xScale={xScale}
            yScale={yScale}
            yMax={yMax}
            minRef={minRef}
            maxRef={maxRef}
            q1Ref={q1Ref}
            q2Ref={q2Ref}
            q3Ref={q3Ref}
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
