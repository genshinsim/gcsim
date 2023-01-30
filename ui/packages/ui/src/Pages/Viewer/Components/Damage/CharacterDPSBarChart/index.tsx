import { Card, FormGroup, HTMLSelect } from "@blueprintjs/core";
import { ElementDPS, FloatStat, SimResults, TargetDPS } from "@gcsim/types";
import { ParentSize } from "@visx/responsive";
import { memo, useState } from "react";
import { CardTitle, NoData, useRefreshWithTimer } from "../../Util";
import { ByCharacterChart, ByCharacterLegend } from "./ByCharacter";
import { ByElementChart, ByElementLegend } from "./ByElement";
import { ByTargetChart, ByTargetLegend } from "./ByTarget";

type GraphData = {
  byElement?: ElementDPS[];
  byCharacter?: FloatStat[];
  byTarget?: TargetDPS[];
}

type Props = {
  data: SimResults | null;
  running: boolean;
  names?: string[];
}

export default ({ data, running, names }: Props) => {
  const [graph, setGraph] = useState("element");
  const [stats, timer] = useRefreshWithTimer(d => {
    return {
      byElement: d?.statistics?.dps_by_element,
      byCharacter: d?.statistics?.character_dps,
      byTarget: d?.statistics?.dps_by_target,
    };
  }, 5000, data, running);

  return (
    <Card className="flex flex-col col-span-full h-96">
      <div className="flex flex-row justify-start gap-5">
        <div className="flex flex-col gap-2">
          <CardTitle title="Character DPS" tooltip="x" timer={timer} />
          <Options graph={graph} setGraph={setGraph} />
        </div>
        <div className="flex flex-grow justify-center items-center">
          <Legend data={stats} names={names} graph={graph} />
        </div>
      </div>
      <ParentSize>
        {({ width, height }) => (
          <Graph data={stats} names={names} width={width} height={height} graph={graph} />
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
  data: GraphData;
  names?: string[];
  graph: string;
  width: number;
  height: number;
}

const Graph = memo(({ data, names, graph, width, height }: GraphProps) => {
  if (graph === "element") {
    return (
      <ByElementChart
          width={width}
          height={height}
          names={names}
          dps={data.byElement} />
    );
  } else if (graph === "character") {
    return (
      <ByCharacterChart
          width={width}
          height={height}
          names={names}
          dps={data.byCharacter} />
    );
  } else if (graph === "target") {
    return (
      <ByTargetChart
          width={width}
          height={height}
          names={names}
          dps={data.byTarget} />
    );
  }
  return <NoData />;
});

type LegendProps = {
  data: GraphData;
  names?: string[];
  graph: string;
}

const Legend = memo(({ data, names, graph }: LegendProps) => {
  if (graph === "element") {
    return <ByElementLegend dps={data.byElement} />;
  } else if (graph === "character") {
    return <ByCharacterLegend names={names} />;
  } else if (graph === "target") {
    return <ByTargetLegend dps={data.byTarget} />;
  }
  return null;
});