import { TargetStats } from "@gcsim/types";
import { LegendOrdinal } from "@visx/legend";
import { scaleOrdinal } from "@visx/scale";
import { useMemo } from "react";
import { DataColors, FloatStatTooltipContent, HorizontalBarStack, NoData } from "../../Util";

type Props = {
  width: number;
  height: number;
  names?: string[];
  dps?: TargetStats[];
}

export const ByTargetLegend = ({ dps }: { dps?: TargetStats[] }) => {
  const keys = useMemo(() => {
    if (dps == null) {
      return [];
    }

    const targets = new Set<string>();
    for (let i = 0; i < dps.length; i++) {
      for (const key in dps[i].targets) {
        targets.add(key);
      }
    }
    return Array.from(targets);
  }, [dps]);

  const scale = scaleOrdinal({
    domain: keys,
    range: keys.map(v => DataColors.target(v)),
  });

  return (
    <LegendOrdinal scale={scale} direction="row" labelMargin="0 15px 0 0" className="flex-wrap" />
  );
};

export const ByTargetChart = ({ width, height, names, dps }: Props) => {
  const { data, keys, xMax } = useData(dps, names);

  if (dps == null || names == null || keys.length == 0) {
    return <NoData />;
  }

  return (
    <HorizontalBarStack<TargetData, string>
        width={width}
        height={height}
        xDomain={[0, xMax]}
        yDomain={names}
        y={d => d.name}
        data={data}
        keys={keys}
        value={(d, k) => {
          if (k in d.data) {
            return d.data[k].mean ?? 0;
          }
          return 0;
        }}
        stat={(d, k) => d.data[k]}
        barColor={k => DataColors.target(k)}
        hoverColor={k => DataColors.targetLabel(k)}
        tooltipContent={(d, k) => (
          <FloatStatTooltipContent
              title={d.name + " target " + k + " dps"}
              data={d.data[k]}
              color={DataColors.targetLabel(k)}
              percent={(d.data[k].mean ?? 0) / d.total}
          />
        )}
    />
  );
};

type TargetData = {
  name: string;
  data: TargetStats;
  total: number;
}

type ChartData = {
  data: TargetData[];
  keys: string[];
  xMax: number;
}

function useData(dps?: TargetStats[], names?: string[]): ChartData {
  return useMemo(() => {
    if (dps == null || names == null) {
      return { data: [], keys: [], xMax: 0 };
    }

    const targets = new Set<string>();
    const data: TargetData[] = [];

    let maxDPS = 0;
    for (let i = 0; i < dps.length; i++) {
      const char = dps[i].targets;
      if (char == null) {
        continue;
      }

      let maxTotal = 0;
      let total = 0;
      for (const key in char) {
        targets.add(key);
        const mean = char[key].mean ?? 0;
        maxTotal += Math.max(char[key].max ?? 0, mean + (char[key].sd ?? 0));
        total += mean;
      }
      maxDPS = Math.max(maxDPS, maxTotal);
      data.push({ name: names[i], data: char, total: total });
    }

    return {
      data: data,
      keys: Array.from(targets),
      xMax: maxDPS,
    };
  }, [dps, names]);
}