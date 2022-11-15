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
import { Card, Colors, HTMLSelect, Icon, NonIdealState, Spinner, SpinnerSize } from "@blueprintjs/core";


export default ({ data }: { data: SimResults | null}) => {
  const [graph, setGraph] = useState(0);

  const titles = [
    <GraphTitle key="dps" title="DPS Distribution" tooltip="test" />,
    <GraphTitle key="rps" title="RPS Distribution" tooltip="test" />,
    <GraphTitle key="eps" title="EPS Distribution" tooltip="test" />,
    <GraphTitle key="hps" title="HPS Distribution" tooltip="test" />,
    <GraphTitle key="sps" title="SPS Distribution" tooltip="test" />,
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
        key="rps"
        data={data?.statistics?.rps}
        color="#147EB3"
        accentColor="#0C5174"
        hoverColor="#68C1EE" />,
    <GraphContent
        key="eps"
        data={data?.statistics?.eps}
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
        key="sps"
        data={data?.statistics?.sps}
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
      <div className="flex flex-row justify-start">
        <HTMLSelect value={graph} onChange={(e) => setGraph(Number(e.target.value))}>
          <option value={0}>DPS</option>
          <option value={2}>EPS</option>
          <option value={1}>RPS</option>
          <option value={3}>HPS</option>
          <option value={4}>SPS</option>
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
    const { i18n } = useTranslation();

    const xMax = width - margin.left - margin.right;
    const yMax = height - margin.top - margin.bottom;
    const numTicks = 7;

    const xScale = useMemo(() => scaleBand<number>({
      range: [0, xMax],
      domain: range(data?.histogram?.length ?? 1),
      paddingInner: 0.05
    }), [data?.histogram?.length, xMax]);

    const yScale = useMemo(() => {
      const max = Math.max(...(data?.histogram ?? [1]));
      return scaleLinear<number>({
        range: [yMax, 0],
        domain: [-3, max + .1 * max],
        // nice: true,
        clamp: true
      });
    }, [data?.histogram, yMax]);

    const meanLine = useMemo(() => {
      if (data?.max == null || data?.min == null || data?.mean == null) {
        return null;
      }

      const xLin = scaleLinear<number>({
        range: [0, xMax],
        domain: [data.min, data.max],
      });
      const x = xLin(data.mean);
      return (
        <g>
          <Line
              from={{ x: x, y: yMax }}
              to={{ x: x, y: 0 }}
              stroke={accentColor}
              strokeWidth={2}
              strokeDasharray="5,2"
              className="opacity-75" />
          <Text x={x} y={yMax} dy="1em" className="fill-gray-400 font-mono" textAnchor="middle">
            {"mean=" + data?.mean?.toLocaleString(i18n.language, { maximumFractionDigits: 2 })}
          </Text>
        </g>
      );
    }, [accentColor, data?.max, data?.mean, data?.min, i18n.language, xMax, yMax]);

    const meanBin = useMemo(() => {
      if (data?.histogram?.length == null || data?.max == null || data?.min == null
          || data?.mean == null) {
        return null;
      }
      const delta = data.histogram.length / (data.max - data.min);
      return Math.floor(delta * (data.mean - data.min));
    }, [data?.histogram?.length, data?.max, data?.mean, data?.min]);

    if (data?.histogram == null) {
      return <NonIdealState icon={<Spinner size={SpinnerSize.LARGE} />} />;
    }

    return (
      <>
        <svg width={width} height={height}>
          <Group left={margin.left} top={margin.top}>
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
            {meanLine}
            {data?.histogram?.map((c, i) => {
              const barWidth = xScale.bandwidth();
              const barHeight = yMax - yScale(c);

              if (c <= 0 || barHeight < 0) {
                return null;
              }

              const barX = xScale(i) ?? 0;
              const barY = yMax - barHeight;
              let fill = color;
              if (i === meanBin) {
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
                    height={barHeight}
                    onMouseLeave={() => hideTooltip()}
                    onMouseMove={() => {
                      if (data?.histogram?.length == null || data?.max == null || data?.min == null) {
                        return null;
                      }

                      const width = (data.max - data.min) / data.histogram.length;
                      const lower = data.min + i * width;
                      const upper = data.min + (i+1) * width;
                      return showTooltip({
                          tooltipData: { index: i, count: c, lower: lower, upper: upper },
                          tooltipLeft: barX + margin.left + (barWidth/2),
                          tooltipTop: barY - 10
                      });
                    }} />
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
                content={<TooltipContent {...tooltipData} />}>
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
    <Text x={props.x} y={props.y} dy="0.25em" textAnchor="end">
      {props.formattedValue}
    </Text>
  );
};

const TooltipContent = (props: TooltipData) => {
  const { i18n } = useTranslation();
  const lower = props.lower?.toLocaleString(i18n.language, { maximumFractionDigits: 2 });
  const upper = props.upper?.toLocaleString(i18n.language, { maximumFractionDigits: 2 });

  return (
    <div className="p-2 font-mono text-xs grid grid-cols-[repeat(2,_min-content)] gap-x-2 justify-center">
      <span className="justify-self-end text-gray-400">lower</span><span>{lower}</span>
      <span className="justify-self-end text-gray-400">upper</span><span>{upper}</span>
      <span className="justify-self-end text-gray-400">itrs</span><span>{props.count}</span>
    </div>
  );
};