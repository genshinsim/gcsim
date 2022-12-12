import { Card, Colors, HTMLSelect } from "@blueprintjs/core";
import { SimResults, SummaryStat } from "@gcsim/types";
import { ParentSize } from "@visx/responsive";
import { useState } from "react";
import { CardTitle } from "../../Util";
import { HistogramGraph } from "./HistogramGraph";

type Props = {
  data: SimResults | null;
}

export default ({ data }: Props) => {
  const [graph, setGraph] = useState("dps");

  return (
    <Card className="col-span-3 min-h-full h-72 min-w-[280px] flex flex-col justify-start gap-2">
      <div className="flex flex-row justify-start">
        <HTMLSelect value={graph} onChange={(e) => setGraph(e.target.value)}>
          <option value="dps">DPS</option>
          <option value="eps">EPS</option>
          <option value="rps">RPS</option>
          <option value="hps">HPS</option>
          <option value="shp">SHP</option>
          <option value="dur">Dur</option>
        </HTMLSelect>
        <div className="flex flex-grow justify-center items-end">
          <GraphTitle graph={graph} />
        </div>
      </div>
      <Graph graph={graph} data={data} />
    </Card>
  );
};

const GraphTitle = ({ graph }: { graph: string }) => {
  if (graph === "dps") {
    return <CardTitle title="DPS Distribution" tooltip="test" />;
  } else if (graph === "eps") {
    return <CardTitle title="EPS Distribution" tooltip="test" />;
  } else if (graph === "rps") {
    return <CardTitle title="RPS Distribution" tooltip="test" />;
  } else if (graph === "hps") {
    return <CardTitle title="HPS Distribution" tooltip="test" />;
  } else if (graph === "shp") {
    return <CardTitle title="SHP Distribution" tooltip="test" />;
  } else if (graph === "dur") {
    return <CardTitle title="Duration Distribution" tooltip="test" />;
  }
  return null;
};

const Graph = ({ graph, data }: { graph: string, data: SimResults | null }) => {
  if (graph === "dps") {
    return (
      <GraphContent
        data={data?.statistics?.dps}
        barColor={Colors.VERMILION3}
        accentColor={Colors.VERMILION1}
        hoverColor={Colors.VERMILION5} />
    );
  } else if (graph === "eps") {
    return (
      <GraphContent
        data={data?.statistics?.eps}
        barColor={Colors.CERULEAN3}
        accentColor={Colors.CERULEAN1}
        hoverColor={Colors.CERULEAN5} />
    );
  } else if (graph === "rps") {
    return (
      <GraphContent
        data={data?.statistics?.rps}
        barColor={Colors.VIOLET3}
        accentColor={Colors.VIOLET1}
        hoverColor={Colors.VIOLET5} />
    );
  } else if (graph === "hps") {
    return (
      <GraphContent
        data={data?.statistics?.hps}
        barColor={Colors.FOREST3}
        accentColor={Colors.FOREST1}
        hoverColor={Colors.FOREST5} />
    );
  } else if (graph === "shp") {
    return (
      <GraphContent
        data={data?.statistics?.shp}
        barColor={Colors.GOLD3}
        accentColor={Colors.GOLD1}
        hoverColor={Colors.GOLD5} />
    );
  } else if (graph === "dur") {
    return (
      <GraphContent
        key="dur"
        data={data?.statistics?.duration}
        barColor={Colors.TURQUOISE3}
        accentColor={Colors.TURQUOISE1}
        hoverColor={Colors.TURQUOISE5} />
    );
  }
  return null;
};

type GraphContentProps = {
  data?: SummaryStat;
  barColor?: string;
  hoverColor?: string;
  accentColor?: string;
}

const GraphContent = (props: GraphContentProps) => {
  return (
    <ParentSize>
      {({ width, height }) => (
        <HistogramGraph
            width={width}
            height={height}
            data={props.data}
            barColor={props.barColor}
            hoverColor={props.hoverColor}
            accentColor={props.accentColor} />
      )}
    </ParentSize>
  );
};