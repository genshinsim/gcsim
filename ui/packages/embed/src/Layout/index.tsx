import { Card } from "@blueprintjs/core";
import { SimResults } from "@gcsim/types";
import { NoDataIcon } from "@gcsim/ui/src/Pages/Viewer/Components/Util/NoData";
import { Avatars } from "../Avatar/AvatarCard";
import { Metadata } from "./Metadata";

type Props = {
  data: SimResults;
};

const Layout = ({ data }: Props) => {
  return (
    <div className="bp4-dark flex flex-col gap-2 p-1 h-screen">
      <Avatars chars={data.character_details} />
      <Metadata data={data} />
      <div className="flex flex-auto flex-row gap-2 justify-end h-full">
        <Card className="w-1/4 p-0 items-center flex justify-center">
          <NoDataIcon className="h-[64px]" />
        </Card>
        <Card className="w-1/4 p-0 items-center flex justify-center">
          <NoDataIcon className="h-[64px]" />
        </Card>
        <Card className="w-1/4 p-0 items-center flex justify-center">
          <NoDataIcon className="h-[64px]" />
        </Card>
        <Card className="w-1/4 p-0 items-center flex justify-center">
          <NoDataIcon className="h-[64px]" />
        </Card>
      </div>
    </div>
  );
};

export default Layout;