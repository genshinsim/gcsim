import { Card, Colors } from "@blueprintjs/core";
import { ElementDPS, FloatStat, SimResults } from "@gcsim/types";
import { ParentSize } from "@visx/responsive";
import { scaleOrdinal } from "@visx/scale";
import { useMemo } from "react";
import { useTranslation } from "react-i18next";
import CardTitle from "../../CardTitle";
import OuterLabelPie from "../../OuterLabelPie";

type ElementColor = {
  label: string;
  value: string;
}

const elements: Map<string, ElementColor> = new Map([
  ["electro", { label: Colors.VIOLET4, value: Colors.VIOLET3 }],
  ["pyro", { label: Colors.VERMILION4, value: Colors.VERMILION3 }],
  ["cryo", { label: "#95CACB", value: "#4B8DAA" }],
  ["hydro", { label: Colors.CERULEAN4, value: Colors.CERULEAN3 }],
  ["dendro",{ label: Colors.FOREST4, value: Colors.FOREST3 }],
  ["anemo",{ label: Colors.TURQUOISE4, value: Colors.TURQUOISE3 }],
  ["geo", { label: Colors.GOLD4, value: Colors.GOLD3 }],
  ["physical",{ label: Colors.SEPIA4, value: Colors.SEPIA3 }],
  // not possible?
  ["frozen", { label: "#000", value: "#000" }],
  ["quicken", { label: "#FFF", value: "#FFF" }],
]);

const color = scaleOrdinal<string, string>({
  domain: Array.from(elements.keys()),
  range: Array.from(elements.values()).map(e => e.value),
});

const labelColor = scaleOrdinal<string, string>({
  domain: Array.from(elements.keys()),
  range: Array.from(elements.values()).map(e => e.label),
});

type Props = {
  data: SimResults | null;
}

export default ({ data }: Props) => {
  return (
    <Card className="flex flex-col col-span-2 h-72 min-h-full gap-0">
      <CardTitle title="Element DPS Breakdown" tooltip="x" />
      <ParentSize>
        {({ width, height }) => (
          <DPSPie width={width} height={height} dps={data?.statistics?.element_dps} />
        )}
      </ParentSize>
    </Card>
  );
};

type PieProps = {
  width: number;
  height: number;
  dps?: ElementDPS;
}

type ElementData = {
  label: string;
  value: FloatStat;
}

const DPSPie = ({ width, height, dps }: PieProps) => {
  const { i18n } = useTranslation();

  const total = useMemo(() => {
    if (dps == null) {
      return 0;
    }

    let out = 0;
    for (const key in dps) {
      out += dps[key].mean ?? 0;
    }
    return out;
  }, [dps]);

  const data: ElementData[] = useMemo(() => {
    if (dps == null) {
      return [];
    }
    const out: ElementData[] = [];
    for (const key in dps) {
      out.push({ label: key, value: dps[key] });
    }
    return out;
  }, [dps]);

  if (dps == null) {
    return null;
  }

  return (
    <OuterLabelPie
        width={width}
        height={height}
        data={data}
        pieValue={d => (d.value.mean ?? 0) / total}
        color={d => color(d.label)}
        labelColor={d => labelColor(d.label)}
        labelText={d => d.label}
        labelValue={d => {
          return ((d.value.mean ?? 0) / total).toLocaleString(
              i18n.language, { maximumFractionDigits: 0, style: "percent" });
        }}
        tooltipContent={d => <TooltipContent data={d} color={labelColor} />}
    />
  );
};

type TooltipProps = {
  data: ElementData;
  color: (x: string) => string;
}

const TooltipContent = ({ data, color }: TooltipProps) => {
  const colorStr = color(data.label);
  return (
    <div className="flex flex-col px-2 py-1 font-mono text-xs">
      <span style={{ color: colorStr }}>{data.label + " dps"}</span>
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