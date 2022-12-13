import { Card, FormGroup, HTMLSelect } from "@blueprintjs/core";
import { SimResults } from "@gcsim/types";
import { ParentSize } from "@visx/responsive";
import { useState } from "react";
import { CardTitle, NoData } from "../../Util";
import { ByCharacterChart, ByCharacterLegend } from "./ByCharacter";
import { ByElementChart, ByElementLegend } from "./ByElement";
import { ByTargetChart, ByTargetLegend } from "./ByTarget";

type Props = {
  data: SimResults | null;
}

export default ({ data }: Props) => {
  const [graph, setGraph] = useState("element");
  const names = data?.character_details?.map(c => c.name);

  return (
    <Card className="flex flex-col col-span-full h-96">
      <div className="flex flex-row justify-start gap-5">
        <div className="flex flex-col gap-2">
          <CardTitle title="Character DPS" tooltip="x" />
          <Options graph={graph} setGraph={setGraph} />
        </div>
        <div className="flex flex-grow justify-center items-center">
          <Legend data={data} names={names} graph={graph} />
        </div>
      </div>
      <ParentSize>
        {({ width, height }) => (
          <Graph data={data} names={names} width={width} height={height} graph={graph} />
        )}
      </ParentSize>
    </Card>
  );
};

const Options = ({ graph, setGraph }: { graph: string, setGraph: (v: string) => void }) => {
  const label = (
    <span className="text-xs font-mono text-gray-400">
      Grouping
    </span>
  );

  return (
    <FormGroup label={label} inline={true} className="!mb-2">
      <HTMLSelect value={graph} onChange={(e) => setGraph(e.target.value)}>
        <option value={"character"}>Character</option>
        <option value={"element"}>Element</option>
        <option value={"target"}>Target</option>
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

const Graph = ({ data, names, graph, width, height }: GraphProps) => {
  if (graph === "element") {
    return (
      <ByElementChart
          width={width}
          height={height}
          names={names}
          dps={data?.statistics?.dps_by_element} />
    );
  } else if (graph === "character") {
    return (
      <ByCharacterChart
          width={width}
          height={height}
          names={names}
          dps={data?.statistics?.character_dps} />
    );
  } else if (graph === "target") {
    return (
      <ByTargetChart
          width={width}
          height={height}
          names={names}
          dps={data?.statistics?.dps_by_target} />
    );
  }
  return <NoData />;
};

type LegendProps = {
  data: SimResults | null;
  names?: string[];
  graph: string;
}

const Legend = ({ data, names, graph }: LegendProps) => {
  if (graph === "element") {
    return <ByElementLegend dps={data?.statistics?.dps_by_element} />;
  } else if (graph === "character") {
    return <ByCharacterLegend names={names} />;
  } else if (graph === "target") {
    return <ByTargetLegend dps={data?.statistics?.dps_by_target} />;
  }
  return null;
};