import { ElementDPS } from "@gcsim/types";
import { LegendOrdinal } from "@visx/legend";
import { scaleOrdinal } from "@visx/scale";
import { useMemo } from "react";
import { DataColors, FloatStatTooltipContent, HorizontalBarStack, NoData } from "../../Util";

type Props = {
  width: number;
  height: number;
  names?: string[];
  dps?: ElementDPS[];
}

export const ByElementLegend = ({ dps }: { dps?: ElementDPS[] }) => {
  const keys = useMemo(() => {
    if (dps == null) {
      return [];
    }

    const elements = new Set<string>();
    for (let i = 0; i < dps.length; i++) {
      for (const key in dps[i]) {
        elements.add(key);
      }
    }
    return Array.from(elements);
  }, [dps]);

  const scale = scaleOrdinal({
    domain: keys,
    range: keys.map(v => DataColors.element(v)),
  });

  return (
    <LegendOrdinal scale={scale} direction="row" labelMargin="0 15px 0 0" className="flex-wrap" />
  );
};

export const ByElementChart = ({ width, height, names, dps }: Props) => {
  const { data, keys, xMax } = useData(dps, names);

  if (dps == null || names == null || keys.length == 0) {
    return <NoData />;
  }

  return (
    <HorizontalBarStack<ElementData, string>
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
        bottomLabel="DPS"
        barColor={DataColors.element}
        hoverColor={DataColors.elementLabel}
        tooltipContent={(d, k) => (
          <FloatStatTooltipContent
              title={d.name + " " + k + " dps"}
              data={d.data[k]}
              color={DataColors.elementLabel(k)}
              percent={(d.data[k].mean ?? 0) / d.total}
          />
        )}
    />
  );
};

type ElementData = {
  name: string;
  data: ElementDPS;
  total: number;
};

type ChartData = {
  data: ElementData[];
  keys: string[];
  xMax: number;
}

function useData(dps?: ElementDPS[], names?: string[]): ChartData {
  return useMemo(() => {
    if (dps == null || names == null) {
      return { data: [], keys: [], xMax: 0 };
    }

    const elements = new Set<string>();
    const data: ElementData[] = [];

    let maxDPS = 0;
    for (let i = 0; i < dps.length; i++) {
      const char = dps[i];
      let maxTotal = 0;
      let total = 0;
      for (const key in char) {
        elements.add(key);
        const mean = char[key].mean ?? 0;
        maxTotal += Math.max(char[key].max ?? 0, mean + (char[key].sd ?? 0));
        total += mean;
      }
      maxDPS = Math.max(maxDPS, maxTotal);
      data.push({ name: names[i], data: char, total: total });
    }

    return {
      data: data,
      keys: Array.from(elements),
      xMax: maxDPS,
    };
  }, [dps, names]);
}