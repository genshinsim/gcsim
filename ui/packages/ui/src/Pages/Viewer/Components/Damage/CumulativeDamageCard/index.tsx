import { Card, FormGroup, HTMLSelect } from "@blueprintjs/core";
import { SimResults } from "@gcsim/types";
import { ParentSize } from "@visx/responsive";
import { useMemo, useState } from "react";
import { CardTitle, useRefreshWithTimer } from "../../Util";
import { CumulativeGraph, CumulativeLegend } from "./CumulativeDamage";

type Props = {
  data: SimResults | null;
  running: boolean;
};

export default ({ data, running }: Props) => {
  const [graph, setGraph] = useState("overall");
  const [stats] = useRefreshWithTimer(
    (d) => ({
      cumu: d?.statistics?.cumu_damage,
    }),
    250,
    data,
    running
  );

  const targets = useMemo(() => {
    if (stats.cumu?.targets == null) {
      return [];
    }

    const targets = new Set<string>();
    Object.keys(stats.cumu?.targets).forEach((key) => targets.add(key));
    return Array.from(targets);
  }, [stats.cumu]);
  const [target, setTarget] = useState("1");

  return (
    <Card className="flex flex-col col-span-full h-[450px]">
      <div className="flex flex-row justify-start gap-5">
        <div className="flex flex-col gap-2">
          <CardTitle title="Cumulative Damage" tooltip="x" />
          <Options
            graph={graph}
            setGraph={setGraph}
            target={target}
            setTarget={setTarget}
            targets={targets}
          />
        </div>
        <div className="flex flex-grow justify-center items-center">
          <CumulativeLegend />
        </div>
      </div>
      <ParentSize>
        {({ width, height }) => (
          <CumulativeGraph
            width={width}
            height={height}
            graph={graph}
            target={target}
            input={stats.cumu}
          />
        )}
      </ParentSize>
    </Card>
  );
};

const Options = ({
  graph,
  setGraph,
  target,
  setTarget,
  targets,
}: {
  graph: string;
  setGraph: (v: string) => void;
  target: string;
  setTarget: (v: string) => void;
  targets: string[];
}) => {
  const graphLabel = (
    <span className="text-xs font-mono text-gray-400">Type</span>
  );
  const targetLabel = (
    <span className="text-xs font-mono text-gray-400">Target</span>
  );

  return (
    <div className="flex flex-row gap-4">
      <FormGroup label={graphLabel} inline={true} className="!mb-2">
        <HTMLSelect value={graph} onChange={(e) => setGraph(e.target.value)}>
          <option value={"overall"}>Overall</option>
          <option value={"target"}>Target</option>
        </HTMLSelect>
      </FormGroup>
      {graph === "target" ? (
        <FormGroup label={targetLabel} inline={true} className="!mb-2">
          <HTMLSelect
            value={target}
            onChange={(e) => setTarget(e.target.value)}
          >
            {targets.map((target) => (
              <option key={target} value={target}>
                {target}
              </option>
            ))}
          </HTMLSelect>
        </FormGroup>
      ) : null}
    </div>
  );
};
