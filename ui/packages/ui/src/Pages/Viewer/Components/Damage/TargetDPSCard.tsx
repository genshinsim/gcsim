import { Card } from "@blueprintjs/core";
import { FloatStat, SimResults, TargetDPS } from "@gcsim/types";
import { ParentSize } from "@visx/responsive";
import { memo, useMemo } from "react";
import { useTranslation } from "react-i18next";
import { CardTitle, useDataColors, FloatStatTooltipContent, NoData, OuterLabelPie, useRefreshWithTimer } from "../Util";

type Props = {
  data: SimResults | null;
  running: boolean;
}

export default ({ data, running }: Props) => {
  const { t } = useTranslation();
  const [dps, timer] = useRefreshWithTimer(
      d => d?.statistics?.target_dps, 10000, data, running);

  return (
    <Card className="flex flex-col col-span-2 h-72 min-h-full gap-0">
      <CardTitle title={t<string>("result.dist", { d: t<string>("result.target_dps") })} tooltip="x" timer={timer}/>
      <DPSPie dps={dps} />
    </Card>
  );
};

type PieProps = {
  dps?: TargetDPS;
}

const DPSPie = memo(({ dps }: PieProps) => {
  const { DataColors } = useDataColors();
  const { i18n, t } = useTranslation();
  const { data } = useData(dps);

  if (dps == null) {
    return <NoData />;
  }

  return (
    <ParentSize>
      {({ width, height }) => (
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
                  title={t<string>("viewer.target") + " " + d.label + " DPS"}
                  data={d.value}
                  color={DataColors.targetLabel(d.label)}
                  percent={d.pct} />
            )}
        />
      )}
    </ParentSize>
  );
});

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