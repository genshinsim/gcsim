import { Colors, Icon } from "@blueprintjs/core";
import { Tooltip2 } from "@blueprintjs/popover2";
import { memo } from "react";

type Props = {
  title: string;
  tooltip?: string | JSX.Element;
}

const CardTitle = ({ title, tooltip }: Props) => {
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

export default memo(CardTitle);