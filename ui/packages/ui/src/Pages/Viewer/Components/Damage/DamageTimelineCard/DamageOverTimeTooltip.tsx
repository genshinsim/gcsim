import { Point } from "./DamageOverTimeData";
import { ScaleLinear } from "d3-scale";
import { localPoint } from "@visx/event";
import { Group } from "@visx/group";
import { Line } from "@visx/shape";
import { DataColors } from "../../Util";
import { MutableRefObject } from "react";
import { useTranslation } from "react-i18next";
import { Colors } from "@blueprintjs/core";
import { TooltipWithBounds } from "@visx/tooltip";
import { Popover2 } from "@blueprintjs/popover2";

export interface TooltipData {
  index: number;
}

export interface TooltipHandles {
  mouseLeave: () => void;
  mouseHover: (e: React.MouseEvent) => void;
  clearTimeout: () => void;
}

export function useTooltipHandles(
      showTooltip: (args: ShowTooltipArgs<TooltipData>) => void,
      hideTooltip: () => void,
      xScale: ScaleLinear<number, number>,
      yMax: number,
      margin: { left: number, right: number, top: number, bottom: number },
      bucketSize: number,
    ): TooltipHandles {
  let tooltipTimeout: number;
  const mouseLeave = () => {
    tooltipTimeout = window.setTimeout(() => {
      hideTooltip();
    }, 150);
  };

  const clearTimeout = () => {
  if (tooltipTimeout) {
    window.clearTimeout(tooltipTimeout);
  }
  };

  const mouseHover = (e: React.MouseEvent) => {
    const { x } = localPoint(e) || { x: 0 };
    const index = Math.round((60 * xScale.invert(x - margin.left )) / bucketSize);

    clearTimeout();
    showTooltip({
      tooltipData: { index: index },
      tooltipLeft: x,
      tooltipTop: e.nativeEvent.offsetY - 50,
    });
  };

  return {
    mouseLeave: mouseLeave,
    mouseHover: mouseHover,
    clearTimeout: clearTimeout,
  };
}

type ShowTooltipArgs<Datum> = {
  tooltipData?: Datum;
  tooltipLeft?: number;
  tooltipTop?: number;
}

type HoverLineProps = {
  data: Point[];
  xScale: ScaleLinear<number, number>;
  yScale: ScaleLinear<number, number>;
  yMax: number;
  minRef: MutableRefObject<SVGPathElement | null>;
  meanRef: MutableRefObject<SVGPathElement | null>;
  maxRef: MutableRefObject<SVGPathElement | null>;
  tooltipData?: TooltipData;
  tooltipOpen?: boolean;
  tooltipLeft?: number;
  margin: { left: number, right: number, top: number, bottom: number };
}

export const HoverLine = (props: HoverLineProps) => {
  if (!props.tooltipOpen || !props.tooltipLeft || !props.tooltipData || !props.meanRef.current || 
      !props.minRef.current || !props.maxRef.current) {
    return null;
  }

  const x = props.tooltipLeft;
  const point = props.data[props.tooltipData.index];

  return (
    <Group left={-props.margin.left}>
      <Line
          from={{ x: x, y: 0 }}
          to={{ x: x, y: props.yMax }}
          stroke="#FFF"
          opacity={0.5}
          strokeWidth={2}
          pointerEvents="none"
          strokeDasharray="5 2"
      />
      {/* <DataPoint
          cx={x}
          x={props.xScale(point.x)}
          fill={DataColors.qualitative3(1)}
          path={props.maxRef}
          name={"dps-max"}
      />
      <DataPoint
          cx={x}
          x={props.xScale(point.x)}
          fill={DataColors.qualitative2(3)}
          path={props.minRef}
          name={"dps-min"}
      />
      <DataPoint
          cx={x}
          x={props.xScale(point.x)}
          fill={DataColors.qualitative3(8)}
          path={props.meanRef}
          name={"dps-mean"}
      /> */}
    </Group>
  );
};

type DataPointProps = {
  cx: number;
  x: number;
  fill: string;
  path: MutableRefObject<SVGPathElement | null>;
  name: string;
}

const DataPoint = (props: DataPointProps) => {
  if (!props.path.current) {
    return null;
  }

  const y = getPathYFromX(props.x, props.path.current);

  return (
    <g>
      <circle
          cx={props.cx}
          cy={y + 1}
          r={4}
          fill="#000"
          fillOpacity={0.1}
          stroke="#000"
          strokeOpacity={0.1}
          strokeWidth={2}
          pointerEvents="none"
      />
      <circle
          cx={props.cx}
          cy={y}
          r={4}
          fill={props.fill}
          pointerEvents="none"
          stroke="#FFF"
          strokeWidth={2}
      />
    </g>
  );
};

function getPathYFromX(x: number, path: SVGPathElement, error?: number): number {

  error = error || 0.01;
  const maxIterations = 10;

  let lengthStart = 0;
  let lengthEnd = path.getTotalLength();
  let point = path.getPointAtLength((lengthEnd + lengthStart) / 2);
  let iterations = 0;

  while (x < point.x - error || x > point.x + error) {
    const midpoint = (lengthStart + lengthEnd) / 2;

    point = path.getPointAtLength(midpoint);

    if (x < point.x) {
      lengthEnd = midpoint;
    } else {
      lengthStart = midpoint;
    }

    iterations += 1;
    if (maxIterations < iterations) {
      break;
    }
  }
  return point.y;
}

type TooltipProps = {
  data: Point[];
  names?: string[];
  tooltipOpen: boolean;
  tooltipData?: TooltipData;
  tooltipTop?: number;
  tooltipLeft?: number;
  handles: TooltipHandles;
  showTooltip: (args: ShowTooltipArgs<TooltipData>) => void;
  margin: { left: number, right: number, top: number, bottom: number };
}

export const RenderTooltip = (props: TooltipProps) => {
  if (!props.tooltipOpen || !props.tooltipData || !props.tooltipLeft || !props.names) {
    return null;
  }

  const point = props.data[props.tooltipData.index];

  const content = (
    <div
        onMouseMove={() => {
          props.handles.clearTimeout();
          props.showTooltip({
            tooltipData: props.tooltipData,
            tooltipLeft: props.tooltipLeft,
            tooltipTop: props.tooltipTop
          });
        }}
        onMouseLeave={() => props.handles.mouseLeave()}>
      <div className="flex flex-col px-2 py-1 font-mono text-xs">
        <ul className="grid grid-cols-[repeat(2,_min-content)] gap-x-2 justify-start">
         <Item color={Colors.SEPIA4} name="time" value={point.x} suffix={"s"} />
          <Item color={DataColors.qualitative4(3)} name="min" value={point.y.min} />
          <Item color={DataColors.qualitative4(1)} name="max" value={point.y.max} />
          <Item color={DataColors.qualitative4(8)} name="mean" value={point.y.mean} />
          <Item color={DataColors.qualitative4(0)} name="std" value={point.y.sd} />
        </ul>
      </div>
    </div>
  );

  const top = props.tooltipTop;
  const left = props.tooltipLeft;

  return (
    <TooltipWithBounds
        style={{ position: "absolute" }}
        // offsetLeft={props.margin.left}
        left={left}
        top={top}>
      <Popover2
          isOpen={true}
          enforceFocus={false}
          autoFocus={false}
          usePortal={false}
          minimal={true}
          placement="top"
          content={content}>
        <div></div>
      </Popover2>
    </TooltipWithBounds>
  );
};

type ItemProps = {
  name: string;
  value?: number;
  color?: string;
  suffix?: string;
}

const Item = ({ name, value, color, suffix }: ItemProps) => {
  const { i18n } = useTranslation();
  const num = value?.toLocaleString(
        i18n.language, { minimumFractionDigits: 2, maximumFractionDigits: 2 });

  return (
    <>
      <span className="text-gray-400 list-item" style={{ color: color }}>{name}</span>
      <span>{num}{suffix}</span>
    </>
  );
};