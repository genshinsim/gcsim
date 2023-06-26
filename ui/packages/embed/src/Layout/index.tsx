import { SimResults } from "@gcsim/types";
import { useState } from "react";
import { Avatars } from "../Avatar/AvatarCard";
import { Graphs } from "../Graphs";
import { Metadata } from "./Metadata";

type Props = {
  data: SimResults;
};

const Layout = ({ data }: Props) => {
  const [ready, setReady] = useState(false);
  const [loaded, setLoaded] = useState(0);

  const handleLoaded = () => {
    if (loaded + 1 == data.character_details?.length) {
      console.log("all loaded");
      setReady(true);
    }
    setLoaded(loaded + 1);
  };

  const disabled = ready ? "" : "disabled";

  return (
    <div id="card" className={`bp4-dark flex flex-col gap-2 p-1 h-screen ${disabled}`}>
      <Avatars chars={data.character_details} handleLoaded={handleLoaded} />
      <Metadata data={data} />
      <Graphs data={data} />
    </div>
  );
};

export default Layout;