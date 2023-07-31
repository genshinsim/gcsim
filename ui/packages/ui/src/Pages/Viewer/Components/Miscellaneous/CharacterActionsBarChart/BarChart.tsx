import { SourceStat, SourceStats } from "@gcsim/types";
import { LegendOrdinal } from "@visx/legend";
import { scaleOrdinal } from "@visx/scale";
import { memo, useMemo } from "react";
import { DataColors, FloatStatTooltipContent, HorizontalBarStack, NoData } from "../../Util";

type Props = {
  width: number;
  height: number;
  names?: string[];
  actions?: SourceStats[];
  actionNames?: string[] | null;
}

export const BarChartLegend = ({ actionNames }: { actionNames?: string[] | null }) => {
  if (actionNames == null) {
    return null;
  }

  const scale = scaleOrdinal({
    domain: actionNames,
    range: actionNames.map(v => DataColors.action(v)),
  });

  return (
    <LegendOrdinal scale={scale} direction="row" labelMargin="0 15px 0 0" className="flex-wrap" />
  );
};

const Graph = ({ height, width, actions, names, actionNames }: Props) => {
  const {data, xMax} = useData(actions, names);

  if (actions == null || names == null || actionNames == null) {
    return <NoData />;
  }

  return (
    <HorizontalBarStack<ActionData, string>
      width={width}
      height={height}
      xDomain={[0, xMax]}
      yDomain={names}
      y={d => d.name}
      data={data}
      keys={actionNames}
      value={(d, k) => {
        if (k in d.data) {
          return d.data[k].mean ?? 0;
        }
        return 0;
      }}
      stat={(d, k) => d.data[k]}
      barColor={DataColors.action}
      hoverColor={DataColors.actionLabel}
      tooltipContent={(d, k) => (
        <FloatStatTooltipContent
            title={d.name + ": " + k}
            data={d.data[k]}
            color={DataColors.actionLabel(k)}
            percent={(d.data[k].mean ?? 0) / d.total}
        />
      )}
    />
  );
};

export const BarChart = memo(Graph);

type ActionData = {
  name: string;
  data: SourceStat;
  total: number;
};

type ChartData = {
  data: ActionData[];
  xMax: number;
}

function useData(actions?: SourceStats[], names?: string[]): ChartData {
  return useMemo(() => {
    if (actions == null || names == null) {
      return { data: [], keys: [], xMax: 0 };
    }

    const data: ActionData[] = [];

    let maxActions = 0;
    for (let i = 0; i < actions.length; i++) {
      const char = actions[i].sources;
      if (char == null) {
        continue;
      }

      let maxTotal = 0;
      let total = 0;
      for (const key in char) {
        const mean = char[key].mean ?? 0;
        maxTotal += Math.max(char[key].max ?? 0, mean + (char[key].sd ?? 0));
        total += mean;
      }
      maxActions = Math.max(maxActions, maxTotal);
      data.push({ name: names[i], data: char, total: total });
    }

    return {
      data: data,
      xMax: maxActions,
    };
  }, [actions, names]);
}