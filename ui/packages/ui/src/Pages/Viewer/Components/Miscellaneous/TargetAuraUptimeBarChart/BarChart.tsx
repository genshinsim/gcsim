import { FloatStat, SourceStats } from "@gcsim/types";
import { useMemo } from "react";
import {
  useDataColors,
  FloatStatTooltipContent,
  HorizontalBarStack,
  NoData,
} from "../../Util";
import { useTranslation } from "react-i18next";

type Props = {
  width: number;
  height: number;
  auras: string[];
  auraUptime?: SourceStats[];
  target: string;
};

export const BarChart = ({
  width,
  height,
  auras,
  auraUptime,
  target,
}: Props) => {
  const { DataColors } = useDataColors();
  const { t } = useTranslation();
  const { data, sources, xMax } = useData(auras, target, auraUptime);

  const sourceNames = sources.map((s) => s.name);

  if (data == null) {
    return <NoData />;
  }

  return (
    <HorizontalBarStack<Row, string>
      width={width}
      height={height}
      xDomain={[0, xMax]}
      yDomain={sourceNames}
      y={(d) => d.source}
      data={data}
      keys={auras}
      value={(d, k) => {
        if (k in d.data) {
          return d.data[k].data.mean ?? 0;
        }
        return 0;
      }}
      stat={(d, k) => d.data[k].data}
      barColor={DataColors.reactableModifier}
      hoverColor={DataColors.reactableModifierLabel}
      margin={{ top: 0, left: width*0.25, right: width*0.10, bottom: 20 }}
      tooltipContent={(d, k) => (
        <FloatStatTooltipContent
          title={k + " " + t<string>("result.uptime")}
          data={d.data[k].data}
          color={DataColors.reactableModifierLabel(k)}
        />
      )}
      bottomLabel={t<string>("result.p_of_total_dur")}
    />
  );
};

type SourceData = {
  name: string;
  data: FloatStat;
  i: number;
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
  auras: string[],
  target: string,
  auraUptime?: SourceStats[]
): ChartData {
  return useMemo(() => {
    if (auraUptime == null) {
      return { data: [], sources: [], xMax: 0 };
    }

    const rows = new Map<string, SourceData[]>();

    const char = auraUptime[target]?.sources;
    if (char == null) {
      return { data: [], sources: [], xMax: 0 };
    }

    for (const key in char) {
      if (char[key].max == 0) {
        continue;
      }

      const entries = rows.get(key) ?? [];
      entries.push({ name: key, data: char[key], i: Number(target) });
      rows.set(key, entries);
    }

    let maxAura = 0;
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
      maxAura = Math.max(maxAura, max);
    });

    const data = Array.from(rows.values()).map((v) => {
      const total = v.reduce((a, b) => {
        return a + (b.data.mean ?? 0);
      }, 0);

      const m: SourceMap = {};
      let source = "";
      v.forEach((v) => {
        source = v.name;
        m[v.name] = {
          name: v.name,
          data: v.data,
          i: v.i,
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
      xMax: maxAura,
    };
  }, [auraUptime, target]);
}
