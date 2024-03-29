import React from "react";
import { Tooltip } from "../../../ui";

export interface TooltipData {
  index: number;
}

type TooltipHandles = {
  mouseLeave: () => void;
  mouseHover: (e: React.MouseEvent, index: number) => void;
  clearTimeout: () => void;
};

type ShowTooltipArgs<Datum> = {
  tooltipData?: Datum;
  tooltipLeft?: number;
  tooltipTop?: number;
};

export function useTooltipHandles(
  showTooltip: (args: ShowTooltipArgs<TooltipData>) => void,
  hideTooltip: () => void
): TooltipHandles {
  let tooltipTimeout: number;
  const mouseLeave = () => {
    tooltipTimeout = window.setTimeout(() => {
      hideTooltip();
    }, 250);
  };

  const clearTimeout = () => {
    if (tooltipTimeout) {
      window.clearTimeout(tooltipTimeout);
    }
  };

  const mouseHover = (e: React.MouseEvent, index: number) => {
    clearTimeout();
    showTooltip({
      tooltipData: { index: index },
      tooltipLeft: e.nativeEvent.offsetX,
      tooltipTop: e.nativeEvent.offsetY - 35,
    });
  };

  return {
    mouseLeave: mouseLeave,
    mouseHover: mouseHover,
    clearTimeout: clearTimeout,
  };
}

type Props<Datum> = {
  data: Datum[];
  content?: (d: Datum) => string | JSX.Element;
  tooltipOpen: boolean;
  tooltipData?: TooltipData;
  tooltipTop?: number;
  tooltipLeft?: number;
  handles: TooltipHandles;
  showTooltip: (args: ShowTooltipArgs<TooltipData>) => void;
};

export const RenderTooltip = <Datum,>(props: Props<Datum>) => {
  if (!props.tooltipOpen || !props.tooltipData || !props.content) {
    return null;
  }

  const content = (
    <div
      onMouseMove={() => {
        props.handles.clearTimeout();
        props.showTooltip({
          tooltipData: props.tooltipData,
          tooltipLeft: props.tooltipLeft,
          tooltipTop: props.tooltipTop,
        });
      }}
      onMouseLeave={() => props.handles.mouseLeave()}
    >
      {props.content(props.data[props.tooltipData.index])}
    </div>
  );

  return (
    <div
      style={{
        top: props.tooltipTop,
        left: props.tooltipLeft,
        position: "absolute",
      }}
    >
      <Tooltip
        isOpen={true}
        enforceFocus={false}
        autoFocus={false}
        usePortal={false}
        minimal={true}
        placement="top"
        content={content}
      >
        <div></div>
      </Tooltip>
    </div>
  );
};
