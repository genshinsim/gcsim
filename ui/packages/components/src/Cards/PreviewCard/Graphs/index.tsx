import { model } from "@gcsim/types";
import { cn } from "../../../lib/utils";
import { Timeline } from "./Timeline";

type Props = {
  data: model.ISimulationResult;
  className?: string;
};

export const Graphs = ({ data, className = "" }: Props) => {
  const cc = cn(
    "flex flex-auto flex-row gap-2 h-full ml-1 mr-1 mb-1 rounded-sm bg-slate-700",
    className
  );
  return (
    <div className={cc}>
      <div className="w-48 !p-0 items-center flex justify-center">
        <Timeline data={data.statistics?.damage_buckets} />
      </div>
      {/* <Card className="w-1/4 !p-0 items-center flex justify-center">
        <CharacterDPSPie dps={data?.statistics?.character_dps} />
      </Card>
      <Card className="w-1/4 !p-0 items-center flex justify-center">
        <ElementDPSPie dps={data.statistics?.element_dps} />
      </Card>
      <Card className="w-48 !p-0 items-center flex justify-center">
        <Histogram data={data.statistics?.dps} />
      </Card> */}
    </div>
  );
};
