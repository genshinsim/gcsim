import React from "react";
import { Card } from "@blueprintjs/core";
import Header from "../Components/Results/Header";
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
        <Header data={data} />
        <div className="m-w-full min-h-full col-span-2 md:col-span-3 h-96 bg-bp4-black p-5">
          characters
        </div>
        <Card className="m-w-full min-h-full col-span-2 h-32">
          <h5 className="text-xl">summary distributions</h5>
        </Card>
        <div className="m-w-full min-h-full col-span-full h-96 bp4-card">
          timeline
        </div>
      </div>
    </div>
  );
};
