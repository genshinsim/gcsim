import { Card } from "@blueprintjs/core";
import { SimResults } from "@gcsim/types";
import { useState } from "react";
import { Avatars } from "../Avatar/AvatarCard";
import Details from "../Details";

type Props = {
  data: SimResults;
};

const Layout = ({ data }: Props) => {
  const [loaded, setLoaded] = useState(0);
  const [ready, setReady] = useState(false);

  const handleLoaded = () => {
    if (loaded + 1 == data.character_details?.length) {
      console.log("all loaded");
      setReady(true);
    }
    setLoaded(loaded + 1);
  };

  return (
    <div className="bp4-dark flex flex-col align-middle justify-center p-1 h-full">
      <div className="flex flex-row gap-2">
        <Avatars chars={data.character_details} handleLoaded={handleLoaded} />
        <Card className="grow">
          test
        </Card>
      </div>
      <Details data={data} />
    </div>
  );
};

export default Layout;