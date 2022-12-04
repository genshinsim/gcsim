import { Card } from "@blueprintjs/core";
import { FloatStat, SimResults } from "@gcsim/types";
import { ParentSize } from "@visx/responsive";
import { useTranslation } from "react-i18next";
import CardTitle from "../../CardTitle";
import OuterLabelPie from "../../OuterLabelPie";

type Props = {
  data: SimResults | null;
}

export default ({ data }: Props) => {
  return (
    <Card className="flex flex-col col-span-2 h-72 min-h-full gap-0">
      <CardTitle title="Target DPS Breakdown" tooltip="x" />
      <ParentSize>
        {({ width, height }) => (
          <DPSPie width={width} height={height} dps={data?.statistics?.target_dps} />
        )}
      </ParentSize>
    </Card>
  );
};

type PieProps = {
  width: number;
  height: number;
  dps?: FloatStat[];
}

type TargetData = {
  index: number;
  value: FloatStat;
}

const DPSPie = ({ width, height, dps }: PieProps) => {
  const { i18n } = useTranslation();

  if (dps == null) {
    return null;
  }

  const total = dps.reduce((p, a) => p + (a.mean ?? 0), 0);
  const data: TargetData[] = dps.map((value, index) => {
    return {
      index: index,
      value: value,
    };
  });

  const color = (i: number) => ["#147EB3", "#29A634", "#D1980B", "#D33D17"][i];
  const labelColor = (i: number) => ["#3FA6DA", "#43BF4D", "#F0B726", "#EB6847"][i];

  return (
    <OuterLabelPie
        width={width}
        height={height}
        data={data}
        pieValue={d => (d.value.mean ?? 0) / total}
        color={d => color(d.index)}
        labelColor={d => labelColor(d.index)}
        labelText={d => d.index + ""}
        labelValue={d => {
          return ((d.value.mean ?? 0) / total).toLocaleString(
              i18n.language, { maximumFractionDigits: 0, style: "percent" });
        }}
        tooltipContent={d => <TooltipContent data={d} color={labelColor} />}
    />
  );
};

type TooltipProps = {
  data: TargetData;
  color: (x: number) => string;
}

const TooltipContent = ({ data, color }: TooltipProps) => {
  const colorStr = color(data.index);

  return (
    <div className="flex flex-col px-2 py-1 font-mono text-xs">
      <span style={{ color: colorStr }}>{ "target " + data.index + " dps"}</span>
      <ul className="list-disc pl-4 grid grid-cols-[repeat(2,_min-content)] gap-x-2 justify-start">
        <TooltipDataItem color={colorStr} name="mean" value={data.value.mean} />
        <TooltipDataItem color={colorStr} name="min" value={data.value.min} />
        <TooltipDataItem color={colorStr} name="max" value={data.value.max} />
        <TooltipDataItem color={colorStr} name="std" value={data.value.sd} />
      </ul>
    </div>
  );
};

const TooltipDataItem = ({ name, value, color }: { name: string, value?: number, color: string }) => {
  const { i18n } = useTranslation();
  const num = value?.toLocaleString(i18n.language, { maximumFractionDigits: 2 });
  return (
    <>
      <span className="text-gray-400 list-item" style={{color: color}}>{name}</span>
      <span>{num}</span>
    </>
  );
};