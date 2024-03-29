import { model } from "@gcsim/types";
import { cn } from "../../../lib/utils";
import { Histogram } from "./Histogram";
import { CharacterDPSPie, ElementDPSPie } from "./Pie";
import { Timeline } from "./Timeline";

type Props = {
  data: model.SimulationResult;
  className?: string;
};

export const Graphs = ({ data, className = "" }: Props) => {
  const cc = cn(
    "flex flex-auto flex-row gap-1 h-full ml-1 mr-1 mb-1",
    className
  );
  return (
    <div className={cc}>
      <div className="w-48 !p-0 items-center flex justify-center rounded-sm bg-slate-700">
        <Timeline data={data.statistics?.damage_buckets} />
      </div>
      <div className="w-1/4 !p-0 items-center flex justify-center rounded-sm bg-slate-700">
        <CharacterDPSPie dps={data?.statistics?.character_dps} />
      </div>
      <div className="w-1/4 !p-0 items-center flex justify-center rounded-sm bg-slate-700">
        <ElementDPSPie dps={data.statistics?.element_dps} />
      </div>
      <div className="w-48 !p-0 items-center flex justify-center rounded-sm bg-slate-700">
        <Histogram data={data.statistics?.dps} />
      </div>
    </div>
  );
};
