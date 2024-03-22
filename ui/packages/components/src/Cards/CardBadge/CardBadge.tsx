import classNames from "classnames";
import { Badge } from "../../common/ui/badge";

type CardBadgeProps = {
  title?: string;
  value?: string;
  valueCase?: string;
  valueSize?: string;
  bright?: boolean;
  bold?: boolean;
  className?: string;
};

export const CardBadge = ({
  title,
  value,
  bold,
  bright,
  valueCase = "uppercase",
  valueSize = "text-xs",
  className,
}: CardBadgeProps) => {
  if (value === null) {
    return null;
  }
  const titleCls = classNames("mr-2 text-xs lowercase text-mono", {
    "text-gray-400": !bright,
  });
  return (
    <Badge className={className}>
      {title != null && <span className={titleCls}>{title}</span>}
      <span
        className={`${
          bold ? "font-bold" : ""
        } ${valueSize} ${valueCase} text-mono`}
      >
        {value}
      </span>
    </Badge>
  );
};
