import { ReloadIcon } from "@radix-ui/react-icons";
import { memo, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";

type Props = {
  title: string;
  timer?: number;
};

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
};

const TitleWithTooltip = ({ title }: TitleProps) => {
  const out = (
    <div className="flex flex-row text-lg text-gray-400 items-center gap-2 outline-0">
      {title}
    </div>
  );

  return out;
};

const TitleWithTooltipMemo = memo(TitleWithTooltip);

const RefreshStatus = ({ timer }: { timer: number }) => {
  const { t } = useTranslation();
  const [time, setTime] = useState(timeRemaining(timer));

  useEffect(() => {
    const interval = setInterval(() => {
      setTime(timeRemaining(timer));
    }, 500);
    return () => clearInterval(interval);
  }, [timer]);

  return (
    <div className="text-gray-400 outline-0 text-xs flex gap-1 cursor-default">
      <ReloadIcon />
      <span>{time + t("result.seconds_short")}</span>
    </div>
  );
};

function timeRemaining(timer: number) {
  return Math.max(0, Math.ceil((timer - Date.now()) / 1000));
}

export default memo(CardTitle);
