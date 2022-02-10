import { CharDetail } from "../DataType";
import Character from "./Character";

type Props = {
  team: CharDetail[];
};

export default function TeamView(props: Props) {
  const chars = props.team.map((c, i) => {
    if (i > 3) return null; //cant be more than 4
    return <Character key={i} char={c} />;
  });

  return (
    <div className="grid xl:grid-cols-4 lg:grid-cols-4 md:grid-cols-2 sm:grid-cols-2 xs:grid-cols-1 gap-2 m-2 rounded-md">
      {chars}
    </div>
  );
}
