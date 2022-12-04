import { Popover2 } from "@blueprintjs/popover2";
import React from "react";

export interface TooltipData<Datum> {
  index: number;
  data: Datum;
}

type TooltipHandles<Datum> = {
  mouseLeave: () => void;
  mouseHover: (e: React.MouseEvent, index: number, data: Datum) => void;
  clearTimeout: () => void;
}

type ShowTooltipArgs<Datum> = {
  tooltipData?: Datum;
  tooltipLeft?: number;
  tooltipTop?: number;
}

export function useTooltipHandles<Datum>(
    showTooltip: (args: ShowTooltipArgs<TooltipData<Datum>>) => void,
    hideTooltip: () => void): TooltipHandles<Datum> {
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

  const mouseHover = (e: React.MouseEvent, index: number, data: Datum) => {
    clearTimeout();
    showTooltip({
      tooltipData: { index: index, data: data },
      tooltipLeft: e.clientX,
      tooltipTop: e.clientY - 35,
    });
  };

  return {
    mouseLeave: mouseLeave,
    mouseHover: mouseHover,
    clearTimeout: clearTimeout,
  };
}

type Props<Datum> = {
  content?: (d: Datum) => string | JSX.Element;
  tooltipOpen: boolean;
  tooltipData?: TooltipData<Datum>;
  tooltipTop?: number;
  tooltipLeft?: number;
  handles: TooltipHandles<Datum>;
  showTooltip: (args: ShowTooltipArgs<TooltipData<Datum>>) => void;
}

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
        onMouseLeave={() => props.handles.mouseLeave()}>
      {props.content(props.tooltipData.data)}
    </div>
  );

  return (
    <div style={{ top: props.tooltipTop, left: props.tooltipLeft, position: "absolute" }}>
      <Popover2
          isOpen={true}
          enforceFocus={false}
          autoFocus={false}
          usePortal={false}
          minimal={true}
          placement="top"
          popoverClassName="w-36"
          content={content}>
        <div></div>
      </Popover2>
    </div>
  );
};