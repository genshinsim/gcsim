import { Card, FormGroup, HTMLSelect } from "@blueprintjs/core";
import { BucketStats, CharacterBucketStats, SimResults } from "@gcsim/types";
import { ParentSize } from "@visx/responsive";
import { memo, useState } from "react";
import { CardTitle, DataColorsConst, useRefreshWithTimer } from "../../Util";
import { CumulativeGraph, CumulativeLegend } from "./CumulativeContribution";
import { DamageOverTimeGraph, DamageOverTimeLegend } from "./DamageOverTime";
import { useTranslation } from "react-i18next";

export type LegendGlyph = {
  label: string;
  fill: string;
  fillOpacity: number;
  stroke: string;
  strokeOpacity: number;
  strokeDashArray?: string;
};

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
  const { t } = useTranslation();
  const [graph, setGraph] = useState("total");
  const [stats] = useRefreshWithTimer(d => {
    return {
      cumu: d?.statistics?.cumu_damage_contrib,
      dps: d?.statistics?.damage_buckets,
    };
  }, 250, data, running);

  const glyphs: LegendGlyph[] = [
    {
      label: "min",
      fill: DataColorsConst.qualitative2(3),
      fillOpacity: 0.5,
      stroke: DataColorsConst.qualitative2(3),
      strokeOpacity: 0,
    },
    {
      label: "mean",
      fill: DataColorsConst.qualitative3(8),
      fillOpacity: 1.0,
      stroke: DataColorsConst.qualitative3(8),
      strokeOpacity: 0,
    },
    {
      label: "std",
      fill: DataColorsConst.qualitative1(0),
      fillOpacity: 0.2,
      stroke: DataColorsConst.qualitative3(0),
      strokeOpacity: 0.5,
      strokeDashArray: "0 5 0",
    },
    {
      label: "max",
      fill: DataColorsConst.qualitative2(1),
      fillOpacity: 0.35,
      stroke: DataColorsConst.qualitative2(1),
      strokeOpacity: 0,
    },
  ];
  const glyphNames = glyphs.map((g) => g.label);

  return (
    <Card className="flex flex-col col-span-full h-[450px]">
      <div className="flex flex-col sm:flex-row justify-start gap-5">
        <div className="flex flex-col gap-2">
          <CardTitle title={t<string>("result.dmg_timeline")} tooltip="x" />
          <Options graph={graph} setGraph={setGraph} />
        </div>
        <div className="flex flex-grow justify-start sm:justify-center pb-5 sm:pb-0 items-center">
          <Legend graph={graph} names={names} glyphNames={glyphNames} glyphs={glyphs} />
        </div>
      </div>
      <Graph graph={graph} data={stats} names={names} glyphNames={glyphNames} />
    </Card>
  );
};

const Options = ({ graph, setGraph }: { graph: string, setGraph: (v: string) => void }) => {
  const { t } = useTranslation();
  const label = (
    <span className="text-xs font-mono text-gray-400">
      {t<string>("result.type")}
    </span>
  );

  return (
    <FormGroup label={label} inline={true} className="!mb-2">
      <HTMLSelect value={graph} onChange={(e) => setGraph(e.target.value)}>
        <option value={"total"}>{t<string>("result.dmg_over_time")}</option>
        <option value={"cumu"}>{t<string>("result.cumu_contrib")}</option>
      </HTMLSelect>
    </FormGroup>
  );
};

type GraphProps = {
  data: GraphData;
  names?: string[];
  glyphNames: string[];
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
              input={props.data.cumu}
          />
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
              names={props.glyphNames}
              input={props.data.dps} 
          />
        )}
      </ParentSize>
    );
  }
  return null;
});

type LegendProps = {
  names?: string[];
  glyphNames: string[];
  glyphs: LegendGlyph[];
  graph: string;
}

const Legend = memo(({ names, glyphNames, glyphs, graph }: LegendProps) => {
  if (graph === "cumu") {
    return <CumulativeLegend names={names} />;
  } else if (graph === "total") {
    return <DamageOverTimeLegend names={glyphNames} glyphs={glyphs} />;
  }
  return null;
});