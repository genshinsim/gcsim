import { Colors } from "@blueprintjs/core";
import { Popover2 } from "@blueprintjs/popover2";
import { localPoint } from "@visx/event";
import { Group } from "@visx/group";
import { Line } from "@visx/shape";
import { TooltipWithBounds } from "@visx/tooltip";
import { ScaleLinear } from "d3-scale";
import { useTranslation } from "react-i18next";
import { DataColors, FloatStatTooltipContent } from "../../Util";
import { CumulativePoint } from "./CumulativeData";

export interface TooltipData {
  index: number;
}

export interface TooltipHandles {
  mouseLeave: () => void;
  mouseHover: (e: React.MouseEvent) => void;
  clearTimeout: () => void;
}

type ShowTooltipArgs<Datum> = {
  tooltipData?: Datum;
  tooltipLeft?: number;
  tooltipTop?: number;
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
      tooltipTop: yMax,
    });
  };
  
  return {
    mouseLeave: mouseLeave,
    mouseHover: mouseHover,
    clearTimeout: clearTimeout,
  };
}

type HoverLineProps = {
  data: CumulativePoint[];
  xScale: ScaleLinear<number, number>;
  yScale: ScaleLinear<number, number>;
  yMax: number;
  tooltipData?: TooltipData;
  tooltipOpen?: boolean;
  tooltipLeft?: number;
  margin: { left: number, right: number, top: number, bottom: number };
}

export const HoverLine = (props: HoverLineProps) => {
  if (!props.tooltipOpen || !props.tooltipLeft || !props.tooltipData) {
    return null;
  }

  const x = props.tooltipLeft;
  const point = props.data[props.tooltipData.index];

  let total = 0;
  const circles = point.y.map((val, char) => {
    total += val.mean ?? 0;
    const y = props.yScale(total);
    return (
      <g key={"circle-" + char}>
        <circle
            cx={x}
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
            cx={x}
            cy={y}
            r={4}
            fill={DataColors.qualitative4(char)}
            pointerEvents="none"
            stroke="#FFF"
            strokeWidth={2}
        />
      </g>
    );
  });

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
      {circles}
    </Group>
  );
};

type TooltipProps = {
  data: CumulativePoint[];
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
  const { i18n } = useTranslation();

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
      <div className="flex flex-row px-2 py-1 font-mono text-xs gap-2 whitespace-nowrap">
        <span style={{ color: Colors.SEPIA4 }}>time: </span>
        <span>{point.x + "s"}</span>
      </div>
      {point.y.slice(0).reverse().map((val, char) => {
        const i = (props.names?.length ?? 0) - char - 1;
        return (
          <FloatStatTooltipContent
              key={"tooltip-" + i}
              title={props.names?.[i] + " contribution"}
              data={val}
              color={DataColors.characterLabel(i)}
              format={s =>
                s?.toLocaleString(i18n.language, { style: "percent", maximumFractionDigits: 2 })
              }
          />
        );
      })}
    </div>
  );

  const top = props.tooltipTop;
  const left = props.tooltipLeft;

  return (
    <TooltipWithBounds
        style={{ position: "absolute" }}
        offsetLeft={props.margin.left + 50}
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