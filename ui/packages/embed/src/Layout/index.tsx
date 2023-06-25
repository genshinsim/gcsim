import { Card } from "@blueprintjs/core";
import { SimResults } from "@gcsim/types";
import { Avatars } from "../Avatar/AvatarCard";
import Details from "../Details";

type Props = {
  data: SimResults;
};

const Layout = ({ data }: Props) => {
  return (
    <div className="bp4-dark flex flex-col align-middle justify-center p-1 h-full">
      <div className="flex flex-row gap-2">
        <Avatars chars={data.character_details} />
        <Card className="grow">
          test
        </Card>
      </div>
      <Details data={data} />
    </div>
  );
};

export default Layout;