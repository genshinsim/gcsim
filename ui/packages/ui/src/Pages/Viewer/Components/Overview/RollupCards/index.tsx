import { Colors } from "@blueprintjs/core";
import { SimResults } from "@gcsim/types";
import { DPSRollupCard } from "./DPSRollupCard";
import { EPSRollupCard } from "./EPSRollupCard";
import { HPSRollupCard } from "./HPSRollupCard";
import { RPSRollupCard } from "./RPSRollupCard";
import { SHPRollupCard } from "./SHPRollupCard";
import { SimDurRollupCard } from "./SimDurRollupCard";

type Props = {
  data: SimResults | null;
};

// TODO: further optimize by pushing refresh up to here & memo each card impl?
export default ({ data }: Props) => {
  return (
    <div className="col-span-full flex flex-row flex-wrap gap-2 justify-center">
      <DPSRollupCard data={data} color={Colors.VERMILION3} />
      <EPSRollupCard data={data} color={Colors.CERULEAN3} />
      <RPSRollupCard data={data} color={Colors.VIOLET3} />
      <HPSRollupCard data={data} color={Colors.FOREST3} />
      <SHPRollupCard data={data} color={Colors.GOLD3} />
      <SimDurRollupCard data={data} color={Colors.TURQUOISE3} />
    </div>
  );
};
