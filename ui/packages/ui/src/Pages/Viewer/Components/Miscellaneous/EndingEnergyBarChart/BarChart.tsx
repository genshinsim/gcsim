import { EndStats, FloatStat, SourceStats } from "@gcsim/types";
import { range } from "lodash-es";
import { memo, useMemo } from "react";
import { useDataColors, FloatStatTooltipContent, HorizontalBarStack, NoData } from "../../Util";

type Props = {
  width: number;
  height: number;
  names?: string[];
  end_stats?: EndStats[];
}

const Graph = ({ width, end_stats, names }: Props) => {
  const { DataColors } = useDataColors();
  const { data, sources, xMax } = useData(end_stats, names);

  if (end_stats == null || names == null) {
    return <NoData />;
  }

  return (
    <HorizontalBarStack<Row, number>
      width={width}
      height={data.length * 40}
      xDomain={[0, xMax]}
      yDomain={sources}
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
      margin={{ top: 0, left: width*0.15, right: width*0.02, bottom: 20 }}
      tooltipContent={(d, k) => (
        <FloatStatTooltipContent
            title={names[k] + ": " + d.source}
            data={d.data[names[k]].data}
            color={DataColors.characterLabel(k)}
        />
      )}
    />
  );
};

export const BarChart = memo(Graph);

type SourceData = {
  name: string;
  data: FloatStat;
  i: number;
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
  sources: string[]; //character names
  xMax: number;
}

function useData(end_stats?: EndStats[], names?: string[]): ChartData {
  return useMemo(() => {
    if (end_stats == null || names == null) {
      return { data: [], sources: [], xMax: 0 };
    }

    //one row of data per character, source key should just be the char name
    let rows : Row[] = [];

    let xMax = 0;
    end_stats.forEach((e, i) => { 
      const data = e.ending_energy
      if (data == null) {
        //shouldn't happen.. hopefully
        return
      }

      xMax = Math.max(data.max ?? 0,  (data.mean ?? 0) + (data.sd ?? 0), xMax)

      let rowData : SourceMap = {};
      rowData[names[i]] = {
        name: names[i],
        data: data,
        i: i
      }

      rows.push({
        source: names[i],
        mean: data.mean ?? 0, //how often does this happeN/
        data:rowData,
      })

    } )

    return {
      data: rows,
      sources: names,
      xMax: xMax,
    };
  }, [end_stats, names]);
}