import { SimResults } from "@gcsim/types";
import { Avatars } from "../Avatar/AvatarCard";
import { Graphs } from "../Graphs";
import { Metadata } from "./Metadata";

type Props = {
  data: SimResults;
};

const Layout = ({ data }: Props) => {
  return (
    <div id="card" className="bp4-dark flex flex-col gap-2 p-1 h-screen">
      <Avatars chars={data.character_details} />
      <Metadata data={data} />
      <Graphs data={data} />
    </div>
  );
};

export default Layout;