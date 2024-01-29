import { Popover2 } from "@blueprintjs/popover2";
import { SummaryStat } from "@gcsim/types";
import { ScaleLinear, ScaleBand } from "d3-scale";
import { useTranslation } from "react-i18next";

export interface TooltipData {
  x: number;
}

export interface TooltipHandles {
  mouseLeave: () => void;
  mouseHover: (e: React.MouseEvent, data?: SummaryStat) => void;
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
      delta: number | null,
      margin: { left: number, right: number, top: number, bottom: number }
    ): TooltipHandles {
  let tooltipTimeout: number;
  const mouseLeave = () => {
    tooltipTimeout = window.setTimeout(() => {
      hideTooltip();
    }, 750);
  };

  const clearTimeout = () => {
    if (tooltipTimeout) {
      window.clearTimeout(tooltipTimeout);
    }
  };

  const mouseHover = (e: React.MouseEvent, data?: SummaryStat) => {
    if (delta == null || data?.min == null || data?.max == null || data.histogram?.length == null) {
      return null;
    }

    if (e.nativeEvent.offsetX <= margin.left) {
      return null;
    }

    clearTimeout();
    showTooltip({
      tooltipData: { x: e.nativeEvent.offsetX }
    });
  };

  return {
    mouseLeave: mouseLeave,
    mouseHover: mouseHover,
    clearTimeout: clearTimeout,
  };
}

type Props = {
  data?: SummaryStat;
  tooltipOpen: boolean;
  tooltipData?: TooltipData;
  tooltipTop?: number;
  tooltipLeft?: number;
  handles: TooltipHandles;
  showTooltip: (args: ShowTooltipArgs<TooltipData>) => void;
  delta: number | null,
  xLin: ScaleLinear<number, number>,
  xScale: ScaleBand<number>,
  yScale: ScaleLinear<number, number>,
  margin: { left: number, right: number, top: number, bottom: number }
}

export const RenderTooltip = (props: Props) => {
  if (!props.tooltipOpen || !props.tooltipData
      || props.delta == null || props.data?.min == null
      || props.data?.max == null || props.data.histogram?.length == null) {
    return null;
  }

  const temp = props.xLin.invert(props.tooltipData.x - props.margin.left);
  const idx = Math.max(Math.floor(props.delta * (temp - props.data.min)), 0);
  const lower = props.delta == 0 ? props.data.min : props.data.min + idx/props.delta;
  const upper = props.delta == 0 ? props.data.max : props.data.min + (idx+1)/props.delta;
  const count = props.data.histogram[idx];

  if (count <= 0 || idx >= props.data.histogram.length) {
    return null;
  }

  const tooltipLeft = (props.xScale(idx) ?? 0) + props.margin.left + (props.xScale.bandwidth()/2);
  const tooltipTop = props.yScale(count) - 10;

  const content = (
    <div
        onMouseMove={() => {
          props.handles.clearTimeout();
          props.showTooltip({ tooltipData: props.tooltipData });
        }}
        onMouseLeave={() => props.handles.mouseLeave()}>
      <TooltipContent idx={idx} lower={lower} upper={upper} count={count} stat={props.data} />
    </div>
  );

  return (
    <div style={{ top: tooltipTop, left: tooltipLeft, position: "absolute" }}>
      <Popover2
          isOpen={true}
          enforceFocus={false}
          autoFocus={false}
          usePortal={false}
          placement="top"
          content={content}>
        <div></div>
      </Popover2>
    </div>
  );
};

type TooltipContentProps = {
  idx: number;
  lower: number;
  upper: number;
  count: number;
  stat?: SummaryStat;
}

const TooltipContent = ({ idx, lower, upper, count, stat }: TooltipContentProps) => {
  const { i18n, t } = useTranslation();

  const lowerVal = lower?.toLocaleString(i18n.language, { maximumFractionDigits: 2 });
  const upperVal = upper?.toLocaleString(i18n.language, { maximumFractionDigits: 2 });
  const countVal = count?.toLocaleString(i18n.language, { maximumFractionDigits: 2 });
  const mean = stat?.mean?.toLocaleString(i18n.language, { maximumFractionDigits: 2 });
  const p25 = stat?.q1?.toLocaleString(i18n.language, { maximumFractionDigits: 2 });
  const p50 = stat?.q2?.toLocaleString(i18n.language, { maximumFractionDigits: 2 });
  const p75 = stat?.q3?.toLocaleString(i18n.language, { maximumFractionDigits: 2 });

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
        className="px-5 py-2 font-mono text-xs grid grid-cols-[repeat(2,_max-content)] gap-x-2 justify-center">
      {muIndex && muIndex == idx && (
        <><span className="justify-self-end text-gray-400">mean</span><span>{mean}</span></>
      ) || null}
      {p25Index && p25Index == idx && (
        <><span className="justify-self-end text-gray-400">p25</span><span>{p25}</span></>
      ) || null}
      {p50Index && p50Index == idx && (
        <><span className="justify-self-end text-gray-400">p50</span><span>{p50}</span></>
      ) || null}
      {p75Index && p75Index == idx && (
        <><span className="justify-self-end text-gray-400">p75</span><span>{p75}</span></>
      ) || null}
      <span className="justify-self-end text-gray-400">{t<string>("result.lower")}</span><span>{lowerVal}</span>
      <span className="justify-self-end text-gray-400">{t<string>("result.upper")}</span><span>{upperVal}</span>
      <span className="justify-self-end text-gray-400">{t<string>("result.iterations_short")}</span><span>{countVal}</span>
    </div>
  );
};