import { CharacterCard, ConsolidateCharStats } from "~src/Components/Character";
import { CharDetail } from "../DataType";
// import Character from "./Character";

type Props = {
  team: CharDetail[];
};

export default function TeamView(props: Props) {
  const teamStats = ConsolidateCharStats(props.team);

  const chars = props.team.map((c, i) => {
    if (i > 3) return null; //cant be more than 4
    // return <Character key={i} char={c} />;

    return (
      <CharacterCard
        key={c.name}
        char={c}
        stats={teamStats.stats[c.name]}
        statsRows={teamStats.maxRows}
        className="basis-full sm:basis-1/2 hd:basis-1/4 pt-2 pr-2"
      />
    );
  });

  return <div className="flex flex-row flex-wrap pl-2">{chars}</div>;
}
