import { Popover2 } from "@blueprintjs/popover2";
import { TooltipWithBounds } from "@visx/tooltip";
import { useTranslation } from "react-i18next";
import { DataColorsConst } from "../../Util";

export interface TooltipData {
  player: boolean;
  index: number;
  x: number;
  y: number;
  r: number;
}

export interface TooltipHandles {
  mouseLeave: () => void;
  mouseHover: (e: React.MouseEvent, data: TooltipData) => void;
  clearTimeout: () => void;
}

export function useTooltipHandles(
      showTooltip: (args: ShowTooltipArgs) => void,
      hideTooltip: () => void,
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

  const mouseHover = (e: React.MouseEvent, data: TooltipData) => {
    clearTimeout();
    showTooltip({
      tooltipData: data,
      tooltipLeft: e.nativeEvent.offsetX,
      tooltipTop: e.nativeEvent.offsetY - 50,
    });
  };

  return {
    mouseLeave: mouseLeave,
    mouseHover: mouseHover,
    clearTimeout: clearTimeout,
  };
}

type ShowTooltipArgs = {
  tooltipData?: TooltipData;
  tooltipLeft?: number;
  tooltipTop?: number;
}

type Props = {
  tooltipOpen: boolean;
  tooltipData?: TooltipData;
  tooltipTop?: number;
  tooltipLeft?: number;
  handles: TooltipHandles;
  showTooltip: (args: ShowTooltipArgs) => void;
}

export const RenderTooltip = (props: Props) => {
  const { t } = useTranslation();
  if (!props.tooltipOpen || !props.tooltipData) {
    return null;
  }

  const data = props.tooltipData;

  const title = data.player ? t<string>("result.player") : `${t<string>("viewer.target")} ${data.index + 1}`;
  const titleColor = data.player ? DataColorsConst.gray : DataColorsConst.qualitative5(data.index);

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
      <div className="flex flex-col px-2 py-1 font-mono text-xs">
        <span className="text-gray-400 whitespace-nowrap" style={{ color: titleColor }}>
          {title}
        </span>
        <ul className="grid grid-cols-[repeat(2,_max-content)] gap-x-2 justify-start">
          <Item name="x" value={data.x ?? 0} />
          <Item name="y" value={data.y ?? 0} />
          <Item name="r" value={data.r ?? 1} />
        </ul>
      </div>
    </div>
  );

  return (
    <TooltipWithBounds
        style={{ position: "absolute" }}
        left={props.tooltipLeft}
        top={props.tooltipTop}>
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

const Item = ({ name, value }: { name: string, value: number }) => {
  const { i18n } = useTranslation();
  return (
    <>
      <span className="text-gray-400 list-item">{name}</span>
      <span>{value.toLocaleString(
        i18n.language, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}</span>
    </>
  );
};