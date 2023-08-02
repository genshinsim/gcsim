import { Intent, Tag } from "@blueprintjs/core";
import classNames from "classnames";

type ItemProps = {
  title?: string;
  value?: string;
  intent?: Intent;
  bright?: boolean;
  valueCase?: string;
  bold?: boolean;
};

export const Item = ({ title, value, intent, bright, bold, valueCase = "uppercase" }: ItemProps) => {
  if (value == null) {
    return null;
  }

  const titleCls = classNames(
    "text-xs lowercase",
    {"text-gray-400": !bright},
  );

  return (
    <Tag large={true} minimal={!bright} intent={intent}>
      <div className="flex flex-row items-center gap-2 font-mono select-none">
        {title != null && (
          <div className={titleCls}>
            {title}
          </div>
        )}
        <div className={`${bold ? "font-bold" : ""} text-sm ${valueCase}`}>{value}</div>
      </div>
    </Tag>
  );
};