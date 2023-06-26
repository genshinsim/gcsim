import { Card } from "@blueprintjs/core";
import { SimResults } from "@gcsim/types";
import { CharacterDPSPie, ElementDPSPie } from "./Pie";
import { Histogram } from "./Histogram";
import { Timeline } from "./Timeline";

type Props = {
  data: SimResults;
}

export const Graphs = ({ data }: Props) => {
  return (
    <div className="flex flex-auto flex-row gap-2 h-full">
      <Card className="w-48 p-0 items-center flex justify-center">
        <Timeline data={data.statistics?.damage_buckets} />
      </Card>
      <Card className="w-1/4 p-0 items-center flex justify-center">
        <CharacterDPSPie dps={data?.statistics?.character_dps} />
      </Card>
      <Card className="w-1/4 p-0 items-center flex justify-center">
        <ElementDPSPie dps={data.statistics?.element_dps} />
      </Card>
      <Card className="w-48 p-0 items-center flex justify-center">
        <Histogram data={data.statistics?.dps} />
      </Card>
    </div>
  );
};