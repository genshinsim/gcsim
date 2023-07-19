import { Card, FormGroup, HTMLSelect } from "@blueprintjs/core";
import { BucketStats, CharacterBucketStats, SimResults } from "@gcsim/types";
import { ParentSize } from "@visx/responsive";
import { memo, useState } from "react";
import { CardTitle, useRefreshWithTimer } from "../../Util";
import { CumulativeGraph, CumulativeLegend } from "./CumulativeContribution";
import { DamageOverTimeGraph, DamageOverTimeLegend } from "./DamageOverTime";

type GraphData = {
  cumu?: CharacterBucketStats;
  dps?: BucketStats;
}

type Props = {
  data: SimResults | null;
  running: boolean;
  names?: string[];
}

export default ({ data, running, names }: Props) => {
  const [graph, setGraph] = useState("total");
  const [stats] = useRefreshWithTimer(d => {
    return {
      cumu: d?.statistics?.cumu_damage_contrib,
      dps: d?.statistics?.damage_buckets,
    };
  }, 250, data, running);

  return (
    <Card className="flex flex-col col-span-full h-[450px]">
      <div className="flex flex-row justify-start gap-5">
        <div className="flex flex-col gap-2">
          <CardTitle title="Damage Timeline" tooltip="x" />
          <Options graph={graph} setGraph={setGraph} />
        </div>
        <div className="flex flex-grow justify-center items-center">
          <Legend graph={graph} names={names} />
        </div>
      </div>
      <Graph graph={graph} data={stats} names={names} />
    </Card>
  );
};

const Options = ({ graph, setGraph }: { graph: string, setGraph: (v: string) => void }) => {
  const label = (
    <span className="text-xs font-mono text-gray-400">
      Type
    </span>
  );

  return (
    <FormGroup label={label} inline={true} className="!mb-2">
      <HTMLSelect value={graph} onChange={(e) => setGraph(e.target.value)}>
        <option value={"total"}>Damage Over Time</option>
        <option value={"cumu"}>Cumulative Contribution</option>
      </HTMLSelect>
    </FormGroup>
  );
};

type GraphProps = {
  data: GraphData;
  names?: string[];
  graph: string;
}

const Graph = memo((props: GraphProps) => {
  if (props.graph === "cumu") {
    return (
      <ParentSize>
        {({ width, height }) => (
          <CumulativeGraph
              width={width}
              height={height}
              names={props.names}
              input={props.data.cumu} />
        )}
      </ParentSize>
    );
  } else if (props.graph === "total") {
    return (
      <ParentSize>
        {({ width, height }) => (
          <DamageOverTimeGraph
              width={width}
              height={height}
              input={props.data.dps} />
        )}
      </ParentSize>
    );
  }
  return null;
});


const Legend = memo(({ names, graph }: { names?: string[], graph: string }) => {
  if (graph === "cumu") {
    return <CumulativeLegend names={names} />;
  } else if (graph === "total") {
    return <DamageOverTimeLegend />;
  }
  return null;
});