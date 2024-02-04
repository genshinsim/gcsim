import { Card, FormGroup, HTMLSelect } from "@blueprintjs/core";
import { SimResults } from "@gcsim/types";
import { ParentSize } from "@visx/responsive";
import { useMemo, useState } from "react";
import { CardTitle, DataColorsConst, useRefreshWithTimer } from "../../Util";
import { CumulativeGraph, CumulativeLegend } from "./CumulativeDamage";
import { useTranslation } from "react-i18next";

export type LegendGlyph = {
  label: string;
  fill: string;
  fillOpacity: number;
  stroke: string;
  strokeOpacity: number;
  strokeDashArray?: string;
};

type Props = {
  data: SimResults | null;
  running: boolean;
};

export default ({ data, running }: Props) => {
  const { t } = useTranslation();
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

  const glyphs: LegendGlyph[] = [
    {
      label: "min",
      fill: DataColorsConst.qualitative2(3),
      fillOpacity: 0.75,
      stroke: DataColorsConst.qualitative2(3),
      strokeOpacity: 0.75,
      strokeDashArray: "0 5 0",
    },
    {
      label: "max",
      fill: DataColorsConst.qualitative2(1),
      fillOpacity: 0.75,
      stroke: DataColorsConst.qualitative2(1),
      strokeOpacity: 0.75,
      strokeDashArray: "0 5 0",
    },
    {
      label: "p25",
      fill: DataColorsConst.qualitative2(4),
      fillOpacity: 1.0,
      stroke: DataColorsConst.qualitative2(4),
      strokeOpacity: 0,
    },
    {
      label: "p50",
      fill: DataColorsConst.qualitative3(8),
      fillOpacity: 1.0,
      stroke: DataColorsConst.qualitative3(8),
      strokeOpacity: 0,
    },
    {
      label: "p75",
      fill: DataColorsConst.qualitative2(5),
      fillOpacity: 1.0,
      stroke: DataColorsConst.qualitative2(5),
      strokeOpacity: 0,
    },
  ];
  const names = glyphs.map((g) => g.label);

  return (
    <Card className="flex flex-col col-span-full h-[450px]">
      <div className="flex flex-col sm:flex-row justify-start gap-5">
        <div className="flex flex-col gap-2">
          <CardTitle title={t<string>("result.cumu_dmg")} tooltip="x" />
          <Options
            graph={graph}
            setGraph={setGraph}
            target={target}
            setTarget={setTarget}
            targets={targets}
          />
        </div>
        <div className="flex flex-grow justify-start sm:justify-center pb-5 sm:pb-0 items-center">
          <CumulativeLegend names={names} glyphs={glyphs} />
        </div>
      </div>
      <ParentSize>
        {({ width, height }) => (
          <CumulativeGraph
            width={width}
            height={height}
            graph={graph}
            target={target}
            names={names}
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
  const { t } = useTranslation();
  const graphLabel = (
    <span className="text-xs font-mono text-gray-400">{t<string>("result.type")}</span>
  );
  const targetLabel = (
    <span className="text-xs font-mono text-gray-400">{t<string>("viewer.target")}</span>
  );

  return (
    <div className="flex flex-row gap-4">
      <FormGroup label={graphLabel} inline={true} className="!mb-2">
        <HTMLSelect value={graph} onChange={(e) => setGraph(e.target.value)}>
          <option value={"overall"}>{t<string>("result.overall")}</option>
          <option value={"target"}>{t<string>("viewer.target")}</option>
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
