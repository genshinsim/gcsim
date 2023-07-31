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
  const [field_time, timer] = useRefreshWithTimer(
      d => d?.statistics?.field_time, 10000, data, running);
  return (
    <Card className="flex flex-col col-span-3 h-72 min-h-full gap-0">
      <CardTitle title="Field Time Distribution" tooltip="x" timer={timer} />
      <FieldTimePie names={names} field_time={field_time} />
    </Card>
  );
};

type PieProps = {
  names?: string[];
  field_time?: FloatStat[];
}

const FieldTimePie = memo(({ names, field_time }: PieProps) => {
  const { i18n } = useTranslation();
  const { data } = useData(field_time, names);

  if (field_time == null || names == null) {
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
                  title={d.name + " field time (s)"}
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

export function useData(field_time?: FloatStat[], names?: string[]): { data: CharacterData[], total: number } {
  const total = useMemo(() => {
    if (field_time == null) {
      return 0;
    }

    return field_time.reduce((p, a) => p + (a.mean ?? 0), 0);
  }, [field_time]);

  const data: CharacterData[] = useMemo(() => {
    if (field_time == null || names == null) {
      return [];
    }

    return field_time.map((value, index) => {
      return {
        name: names[index],
        index: index,
        value: value,
        pct: (value.mean ?? 0) / total,
      };
    });
  }, [field_time, names, total]);

  return {
    data: data,
    total: total,
  };
}