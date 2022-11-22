import Overview from "../Components/Results/Overview";
import { SimResults } from "@gcsim/types";
import TeamHeader from "../Components/Results/TeamHeader";
import DistributionCard from "../Components/Results/DistributionCard";

type Props = {
  data: SimResults | null;
};

export default ({ data }: Props) => {
  return (
    <div className="w-full 2xl:mx-auto 2xl:container">
      <div className="grid overflow-hidden grid-cols-2 md:grid-cols-5 auto-rows-auto gap-2 px-2">
        <TeamHeader data={data} />
        <Overview data={data} />
        <div className="col-span-2 md:col-span-3 h-72 bg-bp4-black p-5">
        </div>
        <DistributionCard data={data} />
      </div>
    </div>
  );
};
