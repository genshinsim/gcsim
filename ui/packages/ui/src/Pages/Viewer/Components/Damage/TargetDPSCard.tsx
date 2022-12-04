import { Card } from "@blueprintjs/core";
import { FloatStat, SimResults } from "@gcsim/types";
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
      <CardTitle title="Target Damage Breakdown (DPS)" tooltip="x" />
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
        color={d => DataColors.qualitative3(d.index)}
        labelColor={d => DataColors.qualitative4(d.index)}
        labelText={d => d.index + ""}
        labelValue={d => {
          return d.pct.toLocaleString(
              i18n.language, { maximumFractionDigits: 0, style: "percent" });
        }}
        tooltipContent={d => (
          <FloatStatTooltipContent
              title={"target " + d.index + " dps"}
              data={d.value}
              color={DataColors.qualitative4(d.index)}
              percent={d.pct} />
        )}
    />
  );
};

type TargetData = {
  index: number;
  value: FloatStat;
  pct: number;
}

function useData(dps?: FloatStat[]): { data: TargetData[], total: number } {
  const total = useMemo(() => {
    if (dps == null) {
      return 0;
    }

    return dps.reduce((p, a) => p + (a.mean ?? 0), 0);
  }, [dps]);

  const data: TargetData[] = useMemo(() => {
    if (dps == null) {
      return [];
    }

    return dps.map((value, index) => {
      return {
        index: index,
        value: value,
        pct: (value.mean ?? 0) / total,
      };
    });
  }, [dps, total]);

  return {
    data: data,
    total: total,
  };
}