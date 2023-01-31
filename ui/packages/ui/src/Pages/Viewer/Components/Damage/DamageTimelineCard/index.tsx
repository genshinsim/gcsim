import { Card, FormGroup, HTMLSelect } from "@blueprintjs/core";
import { FloatStat, SimResults } from "@gcsim/types";
import { ParentSize } from "@visx/responsive";
import { memo, useState } from "react";
import { CardTitle, NoData, useRefreshWithTimer } from "../../Util";
import { CumulativeGraph, CumulativeLegend } from "./CumulativeContribution";

type GraphData = {
  cumu?: FloatStat[][];
}

type Props = {
  data: SimResults | null;
  running: boolean;
  names?: string[];
}

export default ({ data, running, names }: Props) => {
  const [graph, setGraph] = useState("cumu");
  const [stats, timer] = useRefreshWithTimer(d => {
    return {
      cumu: d?.statistics?.cumu_damage_contrib,
    };
  }, 5000, data, running);

  return (
    <Card className="flex flex-col col-span-full h-[450px]">
      <div className="flex flex-row justify-start gap-5">
        <div className="flex flex-col gap-2">
          <CardTitle title="Damage Timeline" tooltip="x" timer={timer} />
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
    return <NoData />;
  }
  return null;
});


const Legend = memo(({ names, graph }: { names?: string[], graph: string }) => {
  if (graph === "cumu") {
    return <CumulativeLegend names={names} />;
  } else if (graph === "total") {
    return null;
  }
  return null;
});