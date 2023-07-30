import { FloatStat, SourceStats } from "@gcsim/types";
import { LegendOrdinal } from "@visx/legend";
import { scaleOrdinal } from "@visx/scale";
import { range } from "lodash-es";
import { memo, useMemo } from "react";
import { DataColors, FloatStatTooltipContent, HorizontalBarStack, NoData } from "../../Util";

type Props = {
  width: number;
  height: number;
  names?: string[];
  energy?: SourceStats[];
}

const margin = { top: 0, left: 300, right: 20, bottom: 40 };

export const BarChartLegend = ({ names }: { names?: string[] }) => {
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

const Graph = ({ width, energy, names }: Props) => {
  const { data, sources, xMax } = useData(energy, names);

  const sourceNames = sources.map(s => s.name);

  if (energy == null || names == null) {
    return <NoData />;
  }

  return (
    <HorizontalBarStack<Row, number>
      width={width}
      height={data.length * 40}
      xDomain={[0, xMax]}
      yDomain={sourceNames}
      y={d => d.source}
      data={data}
      keys={range(names.length)}
      value={(d, k) => {
        if (names[k] in d.data) {
          return d.data[names[k]].data.mean ?? 0;
        }
        return 0;
      }}
      stat={(d, k) => d.data[names[k]].data}
      barColor={k => DataColors.character(k)}
      hoverColor={k => DataColors.characterLabel(k)}
      margin={margin}
      tooltipContent={(d, k) => (
        <FloatStatTooltipContent
            title={names[k] + ": " + d.source}
            data={d.data[names[k]].data}
            color={DataColors.characterLabel(k)}
            percent={d.data[names[k]].pct}
        />
      )}
    />
  );
};

export const BarChart = memo(Graph);

type SourceData = {
  name: string;
  data: FloatStat;
  char: string;
  i: number;
  pct: number;
}

type SourceName = {
  name: string;
  mean: number;
}

type Row = {
  data: SourceMap;
  source: string;
  mean: number;
};

type SourceMap = {
  [key: string]: SourceData;
}

type ChartData = {
  data: Row[];
  sources: SourceName[];
  xMax: number;
}

function useData(energy?: SourceStats[], names?: string[]): ChartData {
  return useMemo(() => {
    if (energy == null || names == null) {
      return { data: [], sources: [], xMax: 0 };
    }

    const rows = new Map<string, SourceData[]>();
    for (let i = 0; i < energy.length; i++) {
      if (names[i] == "") {
        continue;
      }

      const char = energy[i].sources;
      if (char == null) {
        continue;
      }
      
      for (const key in char) {
        if (char[key].max == 0) {
          continue;
        }

        const entries = rows.get(key) ?? [];
        entries.push({ name: key, data: char[key], char: names[i], i: i, pct: 1 });
        rows.set(key, entries);
      }
    }
    
    let maxEnergy = 0;
    const sources: SourceName[] = [];
    rows.forEach((v, k) => {
      const max: number = v.reduce((a, b) => { 
        return a + Math.max(b.data.max ?? 0, (b.data.mean ?? 0) + (b.data.sd ?? 0));
      }, 0);

      const mean = v.reduce((a, b) => {
        return a + (b.data.mean ?? 0);
      }, 0);

      sources.push({ name: k, mean: mean });
      maxEnergy = Math.max(maxEnergy, max);
    });

    const data = Array.from(rows.values())
      .map(v => {
        const total = v.reduce((a, b) => {
          return a + (b.data.mean ?? 0);
        }, 0);

        const m: SourceMap = {};
        let source = "";
        v.forEach(v => {
          source = v.name;
          m[v.char] = {
            name: v.name,
            data: v.data,
            char: v.char,
            i: v.i,
            pct: (v.data.mean ?? 0) / total,
          };
        });

        return {
          data: m,
          source: source,
          mean: total,
        };
      });

    return {
      data: data,
      sources: sources.sort((a, b) => a.mean - b.mean).filter(v => v.mean > 0),
      xMax: maxEnergy
    };
  }, [energy, names]);
}