import { Popover2 } from "@blueprintjs/popover2";

export interface TooltipData<Key> {
  key: Key;
  index: number;
  x: number;
  y: number;
  height: number;
  width: number;
}

export interface TooltipHandles<Key> {
  mouseLeave: () => void;
  mouseHover: (e: React.MouseEvent, data: TooltipData<Key>) => void;
  clearTimeout: () => void;
}

type ShowTooltipArgs<Datum> = {
  tooltipData?: Datum;
  tooltipLeft?: number;
  tooltipTop?: number;
}

export function useTooltipHandles<Key>(
      showTooltip: (args: ShowTooltipArgs<TooltipData<Key>>) => void,
      hideTooltip: () => void
    ): TooltipHandles<Key> {
  let tooltipTimeout: number;
  const mouseLeave = () => {
    tooltipTimeout = window.setTimeout(() => {
      hideTooltip();
    }, 100);
  };

  const clearTimeout = () => {
    if (tooltipTimeout) {
      window.clearTimeout(tooltipTimeout);
    }
  };

  const mouseHover = (e: React.MouseEvent, data: TooltipData<Key>) => {
    clearTimeout();
    showTooltip({
      tooltipData: data,
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

type Props<Datum,Key> = {
  data: Datum[];
  content?: (d: Datum, k: Key) => string | JSX.Element;
  tooltipOpen: boolean;
  tooltipData?: TooltipData<Key>;
  tooltipTop?: number;
  tooltipLeft?: number;
  handles: TooltipHandles<Key>;
  showTooltip: (args: ShowTooltipArgs<TooltipData<Key>>) => void;
}

export const RenderTooltip = <Datum,Key>(props: Props<Datum,Key>) => {
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
      {props.content(props.data[props.tooltipData.index], props.tooltipData.key)}
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
          content={content}>
        <div></div>
      </Popover2>
    </div>
  );
};