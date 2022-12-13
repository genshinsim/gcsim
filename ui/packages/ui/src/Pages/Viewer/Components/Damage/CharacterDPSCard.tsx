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
  const names = data?.character_details?.map(c => c.name);
  return (
    <Card className="flex flex-col col-span-2 h-72 min-h-full gap-0">
      <CardTitle title="Character DPS Distribution" tooltip="x" />
      <ParentSize>
        {({ width, height }) => (
          <DPSPie
              width={width}
              height={height}
              names={names}
              dps={data?.statistics?.character_dps} />
        )}
      </ParentSize>
    </Card>
  );
};

type PieProps = {
  width: number;
  height: number;
  names?: string[];
  dps?: FloatStat[];
}

const DPSPie = ({ width, height, names, dps }: PieProps) => {
  const { i18n } = useTranslation();
  const { data } = useData(dps, names);

  if (dps == null || names == null) {
    return <NoData />;
  }

  return (
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
  );
};

type CharacterData = {
  name: string;
  index: number;
  value: FloatStat;
  pct: number;
}

function useData(dps?: FloatStat[], names?: string[]): { data: CharacterData[], total: number } {
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