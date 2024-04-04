import { model } from "@gcsim/types";
import { useTranslation } from "react-i18next";

type Props = {
  title: string | JSX.Element;
  data: model.DescriptiveStats;
  color?: string;
  percent?: number;
  format?: (n?: number) => string | undefined;
};

export default ({ title, data, color, percent, format }: Props) => {
  return (
    <div className="flex flex-col px-2 py-1 font-mono text-xs">
      <TooltipTitle title={title} color={color} percent={percent} />
      <ul className="list-disc pl-4 grid grid-cols-[repeat(2,_max-content)] gap-x-2 justify-start">
        <Item format={format} color={color} name="mean" value={data.mean} />
        <Item format={format} color={color} name="min" value={data.min} />
        <Item format={format} color={color} name="max" value={data.max} />
        <Item format={format} color={color} name="std" value={data["sd"]} />
      </ul>
    </div>
  );
};

type TitleProps = {
  title: string | JSX.Element;
  color?: string;
  percent?: number;
};

const TooltipTitle = ({ title, color, percent }: TitleProps) => {
  const { i18n } = useTranslation();
  const value = percent?.toLocaleString(i18n.language, {
    maximumFractionDigits: 2,
    style: "percent",
  });

  if (typeof title === "string" || title instanceof String) {
    if (value == null) {
      return (
        <span
          className="text-gray-400 whitespace-nowrap"
          style={{ color: color }}
        >
          {title}
        </span>
      );
    }

    return (
      <div className="flex flex-row flex-nowrap justify-start text-gray-400 gap-2">
        <span className="whitespace-nowrap" style={{ color: color }}>
          {title}
        </span>
        <span>{"(" + value + ")"}</span>
      </div>
    );
  }
  return title;
};

type ItemProps = {
  name: string;
  value?: number | null;
  color?: string;
  format?: (n?: number) => string | undefined;
};

const Item = ({ name, value, color, format }: ItemProps) => {
  const { i18n } = useTranslation();
  let num: string | undefined;
  if (format == null) {
    num = value?.toLocaleString(i18n.language, {
      minimumFractionDigits: 2,
      maximumFractionDigits: 2,
    });
  } else {
    num = format(value ?? 0);
  }

  return (
    <>
      <span className="text-gray-400 list-item" style={{ color: color }}>
        {name}
      </span>
      <span>{num}</span>
    </>
  );
};
