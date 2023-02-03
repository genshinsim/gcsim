import { Colors, Icon } from "@blueprintjs/core";
import { Tooltip2 } from "@blueprintjs/popover2";
import { memo, useEffect, useState } from "react";

type Props = {
  title: string;
  tooltip?: string | JSX.Element;
  timer?: number;
}

const CardTitle = (props: Props) => {
  if (props.timer == null || props.timer == 0) {
    return <TitleWithTooltipMemo {...props} />;
  }

  return (
    <div className="flex flex-row justify-between items-center gap-4">
      <TitleWithTooltipMemo {...props} />
      <RefreshStatus timer={props.timer} />
    </div>
  );
};

type TitleProps = {
  title: string;
  tooltip?: string | JSX.Element;
}

const TitleWithTooltip = ({ title, tooltip }: TitleProps) => {
  const helpIcon = tooltip == null ? null : <Icon icon="help" color={Colors.GRAY1} />;
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

const TitleWithTooltipMemo = memo(TitleWithTooltip);

const RefreshStatus = ({ timer }: { timer: number }) => {
  const [time, setTime] = useState(timeRemaining(timer));
  
  useEffect(() => {
    const interval = setInterval(() => {
      setTime(timeRemaining(timer));
    }, 500);
    return () => clearInterval(interval);
  }, [timer]);

  return (
    <div className="text-gray-400 outline-0 text-xs flex gap-1 cursor-default">
      <Icon icon="refresh" color={Colors.GRAY1} size={12} className="pt-[3px]" />
      <span>{time + "s"}</span>
    </div>
  );
};

function timeRemaining(timer: number) {
  return Math.max(0, Math.ceil((timer - Date.now()) / 1000));
}

export default memo(CardTitle);