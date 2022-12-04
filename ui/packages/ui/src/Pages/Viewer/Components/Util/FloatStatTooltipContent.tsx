import { FloatStat } from "@gcsim/types";
import { useTranslation } from "react-i18next";

type Props = {
  title: string | JSX.Element;
  data: FloatStat;
  color?: string;
  percent?: number;
}

export default ({ title, data, color, percent }: Props) => (
  <div className="flex flex-col px-2 py-1 font-mono text-xs">
    <TooltipTitle title={title} color={color} percent={percent} />
    <ul className="list-disc pl-4 grid grid-cols-[repeat(2,_min-content)] gap-x-2 justify-start">
      <Item color={color} name="mean" value={data.mean} />
      <Item color={color} name="min" value={data.min} />
      <Item color={color} name="max" value={data.max} />
      <Item color={color} name="std" value={data.sd} />
    </ul>
  </div>
);

type TitleProps = {
  title: string | JSX.Element;
  color?: string;
  percent?: number; 
}

const TooltipTitle = ({ title, color, percent }: TitleProps) => {
  const { i18n } = useTranslation();
  const value = percent?.toLocaleString(
      i18n.language, { maximumFractionDigits: 2, style: "percent" });

  if (typeof title === 'string' || title instanceof String) {
    if (value == null) {
      return (
        <span className="text-gray-400" style={{ color: color }}>{title}</span>
      );
    }

    return (
      <div className="flex flex-row justify-start text-gray-400 gap-1">
        <span style={{ color: color }}>
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
  value?: number;
  color?: string;
}

const Item = ({ name, value, color }: ItemProps) => {
  const { i18n } = useTranslation();
  const num = value?.toLocaleString(i18n.language, { maximumFractionDigits: 2 });

  return (
    <>
      <span className="text-gray-400 list-item" style={{ color: color }}>{name}</span>
      <span>{num}</span>
    </>
  );
};