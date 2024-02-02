import { Card, FormGroup, HTMLSelect } from "@blueprintjs/core";
import { SimResults } from "@gcsim/types";
import { ParentSize } from "@visx/responsive";
import { useEffect, useState } from "react";
import { CardTitle, useRefreshWithTimer } from "../../Util";
import { BarChart, BarChartLegend } from "./BarChart";
import { useTranslation } from "react-i18next";

type Props = {
  data: SimResults | null;
  running: boolean;
  names?: string[];
};

type Graphs = Map<string, string>;

export const SourceDPSCard = ({ data, running, names }: Props) => {
  const { t } = useTranslation();
  const graphs: Graphs = new Map([
    ["dps", "DPS"],
    ["damage_instances", t<string>("result.dmg_instances")],
  ]);
  const [graph, setGraph] = useState("dps");

  const all_filter = t<string>("result.all");
  //@ts-ignore
  const filters: string[] = [all_filter, ...(names || [])];
  const [filter, setFilter] = useState(all_filter);

  // account for the possiblity of changing the language after having this card loaded
  // without: on lang switch -> no data -> require changing the char dropdown value to see it again
  useEffect(() => {
    setFilter(all_filter);
  }, [setFilter, all_filter]);

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
      <div className="flex flex-col sm:flex-row justify-start gap-5">
        <div className="flex flex-col gap-2">
          <CardTitle
            title={t<string>("result.source", { s: graphs.get(graph) })}
            tooltip="x"
            timer={timer}
          />
          <div className="flex flex-row gap-4">
            <Options graph={graph} setGraph={setGraph} graphs={graphs} />
            <Filters filter={filter} setFilter={setFilter} filters={filters} />
          </div>
        </div>
        <div className="flex flex-grow justify-start sm:justify-center pb-5 sm:pb-0 items-center">
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
            all_filter={all_filter}
            filter={filter}
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
  const { t } = useTranslation();
  const label = <span className="text-xs font-mono text-gray-400">{t<string>("result.type")}</span>;

  return (
    <FormGroup label={label} inline={true} className="!mb-2">
      <HTMLSelect value={graph} onChange={(e) => setGraph(e.target.value)}>
        {[...graphs.keys()].map((key) => (
          <option key={key} value={key}>
            {graphs.get(key)}
          </option>
        ))}
      </HTMLSelect>
    </FormGroup>
  );
};

const Filters = ({
  filter,
  setFilter,
  filters,
}: {
  filter: string;
  setFilter: (v: string) => void;
  filters: string[];
}) => {
  const { t } = useTranslation();
  const label = (
    <span className="text-xs font-mono text-gray-400">{t<string>("db.character")}</span>
  );

  return (
    <FormGroup label={label} inline={true} className="!mb-2">
      <HTMLSelect value={filter} onChange={(e) => setFilter(e.target.value)}>
        {[...filters].map((key) => (
          <option key={key} value={key}>
            {key}
          </option>
        ))}
      </HTMLSelect>
    </FormGroup>
  );
};
