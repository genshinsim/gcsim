import { Card } from "@blueprintjs/core";
import { ElementDPS, FloatStat, SimResults } from "@gcsim/types";
import { ParentSize } from "@visx/responsive";
import { useMemo } from "react";
import { useTranslation } from "react-i18next";
import { CardTitle, DataColors, FloatStatTooltipContent, NoData, OuterLabelPie } from "../Util";

type Props = {
  data: SimResults | null;
}

export default ({ data }: Props) => {
  return (
    <Card className="flex flex-col col-span-2 h-72 min-h-full gap-0">
      <CardTitle title="Element DPS Distribution" tooltip="x" />
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

const DPSPie = ({ width, height, dps }: PieProps) => {
  const { i18n } = useTranslation();
  const { data } = useData(dps);

  if (dps == null) {
    return <NoData />;
  }

  return (
    <OuterLabelPie
        width={width}
        height={height}
        data={data}
        pieValue={d => d.pct}
        color={d => DataColors.element(d.label)}
        labelColor={d => DataColors.elementLabel(d.label)}
        labelText={d => d.label}
        labelValue={d => {
          return d.pct.toLocaleString(
              i18n.language, { maximumFractionDigits: 0, style: "percent" });
        }}
        tooltipContent={d => (
          <FloatStatTooltipContent
              title={d.label + " dps"}
              data={d.value}
              color={DataColors.elementLabel(d.label)}
              percent={d.pct} />
        )}
    />
  );
};

type ElementData = {
  label: string;
  value: FloatStat;
  pct: number;
}


function useData(dps?: ElementDPS): { data: ElementData[], total: number } {
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
      out.push({
        label: key,
        value: dps[key],
        pct: (dps[key].mean ?? 0) / total,
      });
    }
    return out;
  }, [dps, total]);

  return {
    data,
    total
  };
}