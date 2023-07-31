import { Card, FormGroup, HTMLSelect } from "@blueprintjs/core";
import { SimResults } from "@gcsim/types";
import { ParentSize } from "@visx/responsive";
import { useState } from "react";
import { CardTitle, useRefreshWithTimer } from "../../Util";
import { BarChart, BarChartLegend } from "./BarChart";

type Props = {
  data: SimResults | null;
  running: boolean;
  names?: string[];
};

type Graphs = Map<string, string>;

export const SourceDPSCard = ({ data, running, names }: Props) => {
  const graphs: Graphs = new Map([
    ["dps", "DPS"],
    ["damage_instances", "Damage Instances"],
  ]);
  const [graph, setGraph] = useState("dps");
  const [stats, timer] = useRefreshWithTimer(
    (d) => {
      return {
        dps: d?.statistics?.source_dps,
        damage_instances: d?.statistics?.source_damage_instances,
      };
    },
    5000,
    data,
    running
  );

  const chart_data = graph === "dps" ? stats.dps : stats.damage_instances;

  return (
    <Card className="flex flex-col col-span-full min-h-96">
      <div className="flex flex-row justify-start gap-5">
        <div className="flex flex-col gap-2">
          <CardTitle
            title={`Source ${graphs.get(graph)}`}
            tooltip="x"
            timer={timer}
          />
          <Options graph={graph} setGraph={setGraph} graphs={graphs} />
        </div>
        <div className="flex flex-grow justify-center items-center">
          <BarChartLegend names={names} />
        </div>
      </div>
      <ParentSize>
        {({ width, height }) => (
          <BarChart
            width={width}
            height={height}
            dps={chart_data}
            names={names}
          />
        )}
      </ParentSize>
    </Card>
  );
};

const Options = ({
  graph,
  setGraph,
  graphs,
}: {
  graph: string;
  setGraph: (v: string) => void;
  graphs: Graphs;
}) => {
  const label = <span className="text-xs font-mono text-gray-400">Type</span>;

  return (
    <FormGroup label={label} inline={true} className="!mb-2">
      <HTMLSelect value={graph} onChange={(e) => setGraph(e.target.value)}>
        {[...graphs.keys()].map((key) => (
          <option value={key}>{graphs.get(key)}</option>
        ))}
      </HTMLSelect>
    </FormGroup>
  );
};
