import { FloatStat } from "@gcsim/types";
import { LegendOrdinal } from "@visx/legend";
import { scaleOrdinal } from "@visx/scale";
import { range } from "lodash-es";
import { useMemo } from "react";
import { DataColors, FloatStatTooltipContent, HorizontalBarStack, NoData } from "../../Util";

type Props = {
  width: number;
  height: number;
  names?: string[];
  dps?: FloatStat[];
}

export const ByCharacterLegend = ({ names }: { names?: string[] }) => {
  if (names == null) {
    return null;
  }

  const scale = scaleOrdinal({
    domain: names,
    range: names.map((n, i) => DataColors.character(i)),
  });

  return (
    <LegendOrdinal scale={scale} direction="row" labelMargin="0 15px 0 0" className="flex-wrap" />
  );
};

export const ByCharacterChart = ({ width, height, names, dps }: Props) => {
  const { data, keys, xMax } = useData(dps, names);

  if (dps == null || names == null) {
    return <NoData />;
  }

  return (
    <HorizontalBarStack<CharacterData, number>
        width={width}
        height={height}
        xDomain={[0, xMax]}
        yDomain={names}
        y={d => d.name}
        data={data}
        keys={keys}
        value={(d, k) => {
          if (d.index === k) {
            return d.data.mean ?? 0;
          }
          return 0;
        }}
        stat={d => d.data}
        barColor={k => DataColors.character(k)}
        hoverColor={k => DataColors.characterLabel(k)}
        tooltipContent={(d, k) => (
          <FloatStatTooltipContent
              title={d.name + " dps"}
              data={d.data}
              color={DataColors.characterLabel(k)}
              percent={1}
          />
        )}
    />
  );
};

type CharacterData = {
  name: string;
  data: FloatStat;
  index: number;
}

type ChartData = {
  data: CharacterData[];
  keys: number[];
  xMax: number;
}

function useData(dps?: FloatStat[], names?: string[]): ChartData {
  return useMemo(() => {
    if (dps == null || names == null) {
      return { data: [], keys: [], xMax: 0 };
    }

    let maxDPS = 0;
    const data: CharacterData[] = dps.map((v, i) => {
      const charMax = Math.max(v.max ?? 0, (v.mean ?? 0) + (v.sd ?? 0));
      maxDPS = Math.max(maxDPS, charMax);
      return { name: names[i], data: v, index: i };
    });

    return { data: data, keys: range(names.length), xMax: maxDPS };
  }, [dps, names]);
}