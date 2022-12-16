import { Card, FormGroup, HTMLSelect } from "@blueprintjs/core";
import { SimResults } from "@gcsim/types";
import { ParentSize } from "@visx/responsive";
import { useState } from "react";
import { CardTitle, NoData } from "../../Util";
import { CumulativeGraph, CumulativeLegend } from "./CumulativeContribution";

type Props = {
  data: SimResults | null;
}

export default ({ data }: Props) => {
  const [graph, setGraph] = useState("cumu");
  const names = data?.character_details?.map(c => c.name);

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
      <ParentSize>
        {({ width, height }) => (
          <Graph
              graph={graph}
              width={width}
              height={height}
              data={data}
              names={names} />
        )}
      </ParentSize>
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
  data: SimResults | null;
  names?: string[];
  graph: string;
  width: number;
  height: number;
}

const Graph = (props: GraphProps) => {
  if (props.graph === "cumu") {
    return (
      <CumulativeGraph
          width={props.width}
          height={props.height}
          names={props.names}
          input={props.data?.statistics?.cumu_damage_contrib} />
    );
  } else if (props.graph === "total") {
    return <NoData />;
  }
  return null;
};


const Legend = ({ names, graph }: { names?: string[], graph: string }) => {
  if (graph === "cumu") {
    return <CumulativeLegend names={names} />;
  } else if (graph === "total") {
    return null;
  }
  return null;
};