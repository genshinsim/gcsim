import { Card } from "@blueprintjs/core";
import { FloatStat, SimResults } from "@gcsim/types";
import { ParentSize } from "@visx/responsive";
import { memo, useMemo } from "react";
import { useTranslation } from "react-i18next";
import { CardTitle, DataColors, FloatStatTooltipContent, NoData, OuterLabelPie, useRefreshWithTimer } from "../Util";

type Props = {
  data: SimResults | null;
  running: boolean;
  names?: string[];
}

export default ({ data, running, names }: Props) => {
  const [dps, timer] = useRefreshWithTimer(
      d => d?.statistics?.character_dps, 10000, data, running);

  return (
    <Card className="flex flex-col col-span-2 h-72 min-h-full gap-0">
      <CardTitle title="Character DPS Distribution" tooltip="x" timer={timer} />
      <DPSPie names={names} dps={dps} />
    </Card>
  );
};

type PieProps = {
  names?: string[];
  dps?: FloatStat[];
}

const DPSPie = memo(({ names, dps }: PieProps) => {
  const { i18n } = useTranslation();
  const { data } = useData(dps, names);

  if (dps == null || names == null) {
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
            color={d => DataColors.character(d.index)}
            labelColor={d => DataColors.characterLabel(d.index)}
            labelText={d => d.name}
            labelValue={d => {
              return d.pct.toLocaleString(
                  i18n.language, { maximumFractionDigits: 0, style: "percent" });
            }}
            tooltipContent={d => (
              <FloatStatTooltipContent
                  title={d.name + " dps"}
                  data={d.value}
                  color={DataColors.characterLabel(d.index)}
                  percent={d.pct} />
            )}
        />
      )}
    </ParentSize>
  );
});

type CharacterData = {
  name: string;
  index: number;
  value: FloatStat;
  pct: number;
}

export function useData(dps?: FloatStat[], names?: string[]): { data: CharacterData[], total: number } {
  const total = useMemo(() => {
    if (dps == null) {
      return 0;
    }

    return dps.reduce((p, a) => p + (a.mean ?? 0), 0);
  }, [dps]);

  const data: CharacterData[] = useMemo(() => {
    if (dps == null || names == null) {
      return [];
    }

    return dps.map((value, index) => {
      return {
        name: names[index],
        index: index,
        value: value,
        pct: (value.mean ?? 0) / total,
      };
    });
  }, [dps, names, total]);

  return {
    data: data,
    total: total,
  };
}