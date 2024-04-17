import classNames from "classnames";

import { Badge, BadgeProps } from "../common/ui/badge";
import { cn } from "../lib/utils";

type ItemProps = {
  title?: string;
  value?: string;
  intent?: BadgeProps["variant"];
  bright?: boolean;
  valueCase?: string;
  bold?: boolean;
  className?: string;
};

export const Item = ({
  title,
  value,
  intent = "default",
  bright,
  bold,
  valueCase = "uppercase",
  className = "",
}: ItemProps) => {
  if (value == null) {
    return null;
  }

  const titleCls = classNames("leading-4 align-bottom text-xs lowercase", {
    "text-gray-400": !bright,
  });

  const cc = cn("font-mono", className);

  return (
    <Badge variant={intent} className={cc}>
      <div className="flex flex-row items-center gap-2 select-none">
        {title != null && <div className={titleCls}>{title}</div>}
        <div
          className={`${
            bold ? "font-bold" : ""
          } leading-4 align-bottom text-sm ${valueCase}`}
        >
          {value}
        </div>
      </div>
    </Badge>
  );
};
