import Summary from "../Components/Results/Summary";
import { SimResults } from "@gcsim/types";
import TeamHeader from "../Components/Results/TeamHeader";

type Props = {
  data: SimResults | null;
};

export default ({ data }: Props) => {
  return (
    <div className="w-full 2xl:mx-auto 2xl:container">
      <div className="grid overflow-hidden grid-cols-2 md:grid-cols-5 auto-rows-auto gap-2">
        <TeamHeader data={data} />
        <Summary data={data} />
      </div>
    </div>
  );
};
