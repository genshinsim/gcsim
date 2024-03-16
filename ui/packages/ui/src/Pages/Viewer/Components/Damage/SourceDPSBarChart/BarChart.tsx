import { FloatStat, SourceStats } from "@gcsim/types";
import { LegendOrdinal } from "@visx/legend";
import { scaleOrdinal } from "@visx/scale";
import { range } from "lodash-es";
import { memo, useMemo } from "react";
import {
  useDataColors,
  FloatStatTooltipContent,
  HorizontalBarStack,
  NoData,
} from "../../Util";

type Props = {
  width: number;
  height: number;
  names?: string[];
  all_filter: string;
  filter: string;
  dps?: SourceStats[];
};

export const BarChartLegend = ({ names }: { names?: string[] }) => {
  const { DataColors } = useDataColors();
  if (names == null) {
    return null;
  }

  const scale = scaleOrdinal({
    domain: names,
    range: names.map((n, i) => DataColors.character(i)),
  });

  return (
    <LegendOrdinal
      scale={scale}
      direction="row"
      labelMargin="0 15px 0 0"
      className="flex-wrap"
    />
  );
};

const Graph = ({ width, dps, names, all_filter, filter }: Props) => {
  const { DataColors } = useDataColors();
  const { data, sources, xMax } = useData(all_filter, filter, dps, names);

  let maxChars = 0;
  const sourceNames = sources.map((s) => {
    maxChars = Math.max(maxChars, s.name.length);
    return s.name;
  });

  if (dps == null || names == null) {
    return <NoData />;
  }

  let minwidth = width;
  if (width < 900) {
    minwidth = 900;
  }

  //see we can force this div to be a min width
  return (
    <div style={{ width: minwidth }}>
    <HorizontalBarStack<Row, number>
        width={minwidth}
      height={data.length * 40}
      xDomain={[0, xMax]}
      yDomain={sourceNames}
      y={(d) => d.source}
      data={data}
      keys={range(names.length)}
      value={(d, k) => {
        if (names[k] in d.data) {
          return d.data[names[k]].data.mean ?? 0;
        }
        return 0;
      }}
      stat={(d, k) => d.data[names[k]].data}
      barColor={(k) => DataColors.character(k)}
      hoverColor={(k) => DataColors.characterLabel(k)}
        margin={{
          top: 0,
          left: maxChars * 8,
          right: 10,
          bottom: 20,
        }}
      tooltipContent={(d, k) => (
        <FloatStatTooltipContent
          title={names[k] + ": " + d.source}
          data={d.data[names[k]].data}
          color={DataColors.characterLabel(k)}
          percent={d.data[names[k]].pct}
        />
      )}
    />
    </div>
  );
};

export const BarChart = memo(Graph);

type SourceData = {
  name: string;
  data: FloatStat;
  char: string;
  i: number;
  pct: number;
};

type SourceName = {
  name: string;
  mean: number;
};

type Row = {
  data: SourceMap;
  source: string;
  mean: number;
};

type SourceMap = {
  [key: string]: SourceData;
};

type ChartData = {
  data: Row[];
  sources: SourceName[];
  xMax: number;
};

function useData(
  all_filter: string,
  filter: string,
  dps?: SourceStats[],
  names?: string[]
): ChartData {
  return useMemo(() => {
    if (dps == null || names == null) {
      return { data: [], sources: [], xMax: 0 };
    }

    const rows = new Map<string, SourceData[]>();
    for (let i = 0; i < dps.length; i++) {
      if (names[i] == "") {
        continue;
      }
      if (filter !== all_filter && names[i] !== filter) {
        continue;
      }

      const char = dps[i].sources;
      if (char == null) {
        continue;
      }

      for (const key in char) {
        if (char[key].max == 0) {
          continue;
        }

        const entries = rows.get(key) ?? [];
        entries.push({
          name: key,
          data: char[key],
          char: names[i],
          i: i,
          pct: 1,
        });
        rows.set(key, entries);
      }
    }

    let maxDPS = 0;
    const sources: SourceName[] = [];
    rows.forEach((v, k) => {
      const max: number = v.reduce((a, b) => {
        return (
          a + Math.max(b.data.max ?? 0, (b.data.mean ?? 0) + (b.data.sd ?? 0))
        );
      }, 0);

      const mean = v.reduce((a, b) => {
        return a + (b.data.mean ?? 0);
      }, 0);

      sources.push({ name: k, mean: mean });
      maxDPS = Math.max(maxDPS, max);
    });

    const data = Array.from(rows.values()).map((v) => {
      const total = v.reduce((a, b) => {
        return a + (b.data.mean ?? 0);
      }, 0);

      const m: SourceMap = {};
      let source = "";
      v.forEach((v) => {
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
      sources: sources
        .sort((a, b) => a.mean - b.mean)
        .filter((v) => v.mean > 0),
      xMax: maxDPS,
    };
  }, [all_filter, filter, dps, names]);
}
