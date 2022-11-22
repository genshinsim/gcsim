import { SimResults, SummaryStat } from "@gcsim/types";
import { Group } from "@visx/group";
import { Line } from "@visx/shape";
import { ParentSize } from "@visx/responsive";
import { scaleBand, scaleLinear } from "@visx/scale";
import { withTooltip } from "@visx/tooltip";
import { WithTooltipProvidedProps } from "@visx/tooltip/lib/enhancers/withTooltip";
import { range } from "lodash-es";
import { useMemo, useState } from "react";
import { GridRows } from "@visx/grid";
import { AxisLeft, TickRendererProps } from "@visx/axis";
import { Text } from "@visx/text";
import { Popover2, Tooltip2 } from "@blueprintjs/popover2";
import { useTranslation } from "react-i18next";
import { Card, Colors, HTMLSelect, Icon, NonIdealState, Tab, Tabs } from "@blueprintjs/core";
import { BoxPlot } from "@visx/stats";

export default ({ data }: { data: SimResults | null}) => {
  const [graph, setGraph] = useState(0);

  const titles = [
    <GraphTitle key="dps" title="DPS Distribution" tooltip="test" />,
    <GraphTitle key="eps" title="EPS Distribution" tooltip="test" />,
    <GraphTitle key="rps" title="RPS Distribution" tooltip="test" />,
    <GraphTitle key="hps" title="HPS Distribution" tooltip="test" />,
    <GraphTitle key="shp" title="SHP Distribution" tooltip="test" />,
    <GraphTitle key="dur" title="Duration Distribution" tooltip="test" />,
  ];

  const graphs = [
    <GraphContent
        key="dps"
        data={data?.statistics?.dps}
        color="#D33D17"
        accentColor="#96290D"
        hoverColor="#FF9980" />,
    <GraphContent
        key="eps"
        data={data?.statistics?.eps}
        color="#147EB3"
        accentColor="#0C5174"
        hoverColor="#68C1EE" />,
    <GraphContent
        key="rps"
        data={data?.statistics?.rps}
        color="#9D3F9D"
        accentColor="#5C255C"
        hoverColor="#D69FD6" />,
    <GraphContent
        key="hps"
        data={data?.statistics?.hps}
        color="#29A634"
        accentColor="#1D7324"
        hoverColor="#62D96B" />,
    <GraphContent
        key="shp"
        data={data?.statistics?.shp}
        color="#D1980B"
        accentColor="#5C4405"
        hoverColor="#FBD065" />,
    <GraphContent
        key="dur"
        data={data?.statistics?.duration}
        color="#00A396"
        accentColor="#004D46"
        hoverColor="#7AE1D8" />,
  ];

  return (
    <Card className="col-span-2 min-h-full h-72 flex flex-col justify-start gap-2">
      <Tabs selectedTabId={"data"} className="-mt-3" >
        <Tab id="data" className="focus:outline-none" title="Distributions" />
        <Tab id="meta" className="focus:outline-none" title="Metadata" />
      </Tabs>
      <div className="flex flex-row justify-start">
        <HTMLSelect value={graph} onChange={(e) => setGraph(Number(e.target.value))}>
          <option value={0}>DPS</option>
          <option value={1}>EPS</option>
          <option value={2}>RPS</option>
          <option value={3}>HPS</option>
          <option value={4}>SHP</option>
          <option value={5}>Dur</option>
        </HTMLSelect>
        <div className="flex flex-grow justify-center items-end">
          {titles[graph]}
        </div>
      </div>
      {graphs[graph]}
    </Card>
  );
};

type GraphContentProps = {
  data?: SummaryStat;
  color?: string;
  hoverColor?: string;
  accentColor?: string;
}

const GraphContent = (props: GraphContentProps) => {
  return (
    <ParentSize>
      {({ width, height }) => (
        <HistogramGraph
            width={width}
            height={height}
            data={props.data}
            color={props.color}
            hoverColor={props.hoverColor}
            accentColor={props.accentColor} />
      )}
    </ParentSize>
  );
};

const GraphTitle = ({ title, tooltip }: { title: string, tooltip?: string | JSX.Element }) => {
  const helpIcon = tooltip == undefined ? null : <Icon icon="help" color={Colors.GRAY1} />;
  const out = (
    <div className="flex flex-row text-lg text-gray-400 items-center gap-2 outline-0">
      {title}
      {helpIcon}
    </div>
  );

  if (tooltip != null) {
    return (
      <div onClick={(e) => e.stopPropagation()} className="cursor-pointer">
        <Tooltip2 content={tooltip}>{out}</Tooltip2>
      </div>
    );
  }
  return out;
};

type HistogramProps = {
  width: number;
  height: number;
  color?: string;
  hoverColor?: string;
  accentColor?: string;
  data?: SummaryStat;
  margin?: { left: number, right: number, top: number, bottom: number };
}

type TooltipData = {
  index: number;
  count: number;
  lower: number;
  upper: number;
}

const defaultMargin = { left: 52, right: 1, top: 5, bottom: 16 };

const HistogramGraph = withTooltip<HistogramProps, TooltipData>(
  ({
    width,
    height,
    data,
    color = "#946638",
    hoverColor = "#D0B090",
    accentColor = "#5E4123",
    margin = defaultMargin,
    tooltipOpen,
    tooltipLeft,
    tooltipTop,
    tooltipData,
    hideTooltip,
    showTooltip,
  }: HistogramProps & WithTooltipProvidedProps<TooltipData>) => {
    const xMax = width - margin.left - margin.right;
    const yMax = height - margin.top - margin.bottom;
    const numTicks = 7;
    
    const { xScale, yScale, xLin, delta } = useScales(data, xMax, yMax);
    
    let tooltipTimeout: number;
    const mouseLeaveHandle = () => {
      tooltipTimeout = window.setTimeout(() => {
        hideTooltip();
      }, 750);
    };

    const mouseHoverHandle = (e: React.MouseEvent) => {
      if (delta == null || data?.min == null || data?.max == null || data.histogram?.length == null) {
        return null;
      }

      if (e.nativeEvent.offsetX <= margin.left) {
        return null;
      }

      const temp = xLin.invert(e.nativeEvent.offsetX - margin.left);
      const idx = Math.max(Math.floor(delta * (temp - data.min)), 0);
      const count = data.histogram[idx];
      const lower = delta == 0 ? data.min : data.min + idx/delta;
      const upper = delta == 0 ? data.max : data.min + (idx+1)/delta;

      if (count <= 0) {
        return null;
      }

      if (tooltipTimeout) {
        window.clearTimeout(tooltipTimeout);
      }

      return showTooltip({
        tooltipData: { index: idx, count: count, lower: lower, upper: upper },
        tooltipLeft: (xScale(idx) ?? 0) + margin.left + (xScale.bandwidth()/2),
        tooltipTop: yScale(count) - 10
      });
    };

    if (data?.histogram == null || delta == null) {
      return <NonIdealState icon="pulse" title="Data not found" />;
    }

    return (
      <>
        <svg
            width={width}
            height={height}
            onMouseMove={mouseHoverHandle}
            onMouseLeave={mouseLeaveHandle}>
          <Group
              left={margin.left}
              top={margin.top}>
            <GridRows
                scale={yScale}
                numTicks={numTicks}
                lineStyle={{ opacity: 0.5 }}
                width={xMax}
                stroke="#e0e0e0"
                height={height} />
            <AxisLeft
                hideAxisLine
                hideTicks
                scale={yScale}
                numTicks={numTicks}
                labelOffset={30}
                labelClassName="fill-gray-400 text-lg"
                tickClassName="fill-gray-400 font-mono"
                tickComponent={(props) => <TickLabel {...props} />}
                label="# iterations" />
            <VerticalLine
                x={data?.mean}
                xScale={xLin}
                yMax={yMax}
                color={accentColor}
                className="opacity-75 fill-gray-400 font-mono" />
            <BoxPlot
                valueScale={xLin}
                min={data.min}
                max={data.max}
                firstQuartile={data.q1}
                median={data.q2}
                thirdQuartile={data.q3}
                horizontal={true}
                boxWidth={10}
                top={yMax + 5}
                fill={hoverColor}
                fillOpacity={0.1}
                stroke={hoverColor}
                strokeWidth={1}
                medianProps={{ style: { stroke: hoverColor } }}
            />
            {data?.histogram?.map((c, i) => {
              const barWidth = xScale.bandwidth();
              const barHeight = yMax - yScale(c);
              const barX = xScale(i) ?? 0;
              const barY = yMax - barHeight;

              if (c <= 0 || barHeight < 0 || data.mean == null || data.min == null) {
                return null;
              }
              
              let fill = color;
              if (i === Math.floor(delta * (data.mean - data.min))) {
                fill = accentColor;
              }
              if (i === tooltipData?.index) {
                fill = hoverColor;
              }

              return (
                <rect
                    key={"bin-" + i}
                    fill={fill}
                    x={barX}
                    y={barY}
                    width={barWidth}
                    height={barHeight} />
              );
            })}
          </Group>
        </svg>
        {tooltipOpen && tooltipData && (
          <div style={{ top: tooltipTop, left: tooltipLeft, position: "absolute" }}>
            <Popover2
                isOpen={true}
                enforceFocus={false}
                autoFocus={false}
                usePortal={false}
                placement="top"
                popoverClassName="w-36"
                content={
                  <TooltipContent
                      data={tooltipData}
                      stat={data}
                      showTooltip={() => {
                        if (tooltipTimeout) {
                          clearTimeout(tooltipTimeout);
                        }
                        
                        showTooltip({
                          tooltipData: tooltipData,
                          tooltipLeft: tooltipLeft,
                          tooltipTop: tooltipTop
                        });
                      }} />
                }>
              <div></div>
            </Popover2>
          </div>
        )}
      </>
    );
  }
);

const TickLabel = (props: TickRendererProps) => {
  return (
    <Text x={props.x} y={props.y} dy="0.25em" textAnchor="end" className="cursor-default">
      {props.formattedValue}
    </Text>
  );
};

type TooltipContentProps = {
  data: TooltipData;
  stat?: SummaryStat;
  showTooltip: () => void;
}

const TooltipContent = ({ data, stat, showTooltip}: TooltipContentProps) => {
  const { i18n } = useTranslation();
  const lower = data.lower?.toLocaleString(i18n.language, { maximumFractionDigits: 2 });
  const upper = data.upper?.toLocaleString(i18n.language, { maximumFractionDigits: 2 });
  const mean = stat?.mean?.toLocaleString(i18n.language, { maximumFractionDigits: 2 });
  const p25 = stat?.q1?.toLocaleString(i18n.language, { maximumFractionDigits: 2 });
  const p50 = stat?.q2?.toLocaleString(i18n.language, { maximumFractionDigits: 2 });
  const p75 = stat?.q3?.toLocaleString(i18n.language, { maximumFractionDigits: 2 });

  const mouseHoverHandle = () => {
    showTooltip();
  };

  if (stat?.histogram?.length == null || stat?.max == null || stat?.min == null) {
    return null;
  }

  const delta = stat.histogram.length <= 1 ? 0 : stat.histogram.length / (stat.max - stat.min);
  const muIndex = stat.mean == null ? null : Math.floor(delta * (stat.mean - stat.min));
  const p25Index = stat.q1 == null ? null : Math.floor(delta * (stat.q1 - stat.min));
  const p50Index = stat.q2 == null ? null : Math.floor(delta * (stat.q2 - stat.min));
  const p75Index = stat.q3 == null ? null : Math.floor(delta * (stat.q3 - stat.min));

  return (
    <div
        className="p-2 font-mono text-xs grid grid-cols-[repeat(2,_min-content)] gap-x-2 justify-center"
        onMouseMove={() => mouseHoverHandle()}>
      {muIndex && muIndex == data.index && (
        <><span className="justify-self-end text-gray-400">mean</span><span>{mean}</span></>
      ) || null}
      {p25Index && p25Index == data.index && (
        <><span className="justify-self-end text-gray-400">p25</span><span>{p25}</span></>
      ) || null}
      {p50Index && p50Index == data.index && (
        <><span className="justify-self-end text-gray-400">p50</span><span>{p50}</span></>
      ) || null}
      {p75Index && p75Index == data.index && (
        <><span className="justify-self-end text-gray-400">p75</span><span>{p75}</span></>
      ) || null}
      <span className="justify-self-end text-gray-400">lower</span><span>{lower}</span>
      <span className="justify-self-end text-gray-400">upper</span><span>{upper}</span>
      <span className="justify-self-end text-gray-400">itrs</span><span>{data.count}</span>
    </div>
  );
};

function useScales(data: SummaryStat | undefined, xMax: number, yMax: number) {
  const xScale = useMemo(() => scaleBand<number>({
    range: [0, xMax],
    domain: range(data?.histogram?.length ?? 1),
    paddingInner: 0.05,
  }), [data?.histogram?.length, xMax]);

  const xLin = useMemo(() => scaleLinear<number>({
    range: [0, xMax],
    domain: [data?.min ?? 0, data?.max ?? 1],
  }), [data?.max, data?.min, xMax]);

  const yScale = useMemo(() => {
    const max = Math.max(...(data?.histogram ?? [1000]));
    return scaleLinear<number>({
      range: [yMax, 0],
      domain: [-3, max + .1 * max],
      clamp: true
    });
  }, [data?.histogram, yMax]);

  const delta = useMemo(() => {
    if (data?.histogram?.length == null || data?.max == null || data?.min == null) {
      return null;
    }

    if (data.histogram.length <= 1) {
      return 0;
    }

    return data.histogram.length / (data.max - data.min);
  }, [data?.histogram?.length, data?.max, data?.min]);

  return { xScale: xScale, yScale: yScale, xLin: xLin, delta: delta };
}

type VerticalLineProps = {
  x?: number;
  xScale: (x: number) => number;
  yMax: number;
  color: string;
  label?: string;
  className?: string;
}

const VerticalLine = ({ x, xScale, yMax, color, label, className }: VerticalLineProps) => {
  if (x == null) {
    return null;
  }
  
  const localX = xScale(x);
  const Label = ({}) => {
    if (label == null) return null;
    return (
      <Text x={xScale(localX)} y={yMax} dy="1em" className={className} textAnchor="middle">
        {label}
      </Text>
    );
  };
  
  return (
    <g>
      <Line
          from={{ x: localX, y: yMax }}
          to={{ x: localX, y: 0 }}
          stroke={color}
          strokeWidth={2}
          strokeDasharray="5,2"
          className={className} />
      <Label />
    </g>
  );
};