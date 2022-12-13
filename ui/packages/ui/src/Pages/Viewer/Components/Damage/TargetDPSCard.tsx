import { Card } from "@blueprintjs/core";
import { FloatStat, SimResults, TargetDPS } from "@gcsim/types";
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
      <CardTitle title="Target DPS Distribution" tooltip="x" />
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
  dps?: TargetDPS;
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
        color={d => DataColors.target(d.label)}
        labelColor={d => DataColors.targetLabel(d.label)}
        labelText={d => d.label}
        labelValue={d => {
          return d.pct.toLocaleString(
              i18n.language, { maximumFractionDigits: 0, style: "percent" });
        }}
        tooltipContent={d => (
          <FloatStatTooltipContent
              title={"target " + d.label + " dps"}
              data={d.value}
              color={DataColors.targetLabel(d.label)}
              percent={d.pct} />
        )}
    />
  );
};

type TargetData = {
  label: string;
  value: FloatStat;
  pct: number;
}

function useData(dps?: TargetDPS): { data: TargetData[], total: number } {
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

  const data: TargetData[] = useMemo(() => {
    if (dps == null) {
      return [];
    }

    const out: TargetData[] = [];
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
    data: data,
    total: total,
  };
}