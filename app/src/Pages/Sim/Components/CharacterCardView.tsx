import {
  CharDetail,
  ConsolidateCharStats,
  CharacterCard,
} from "~src/Components/Character";

type Props = {
  chars: CharDetail[];
};

export function CharacterCardView(props: Props) {
  const teamStats = ConsolidateCharStats(props.chars);

  const rows = props.chars.map((c) => {
    return (
      <CharacterCard
        key={c.name}
        char={c}
        stats={teamStats.stats[c.name]}
        statsRows={teamStats.maxRows}
        className="basis-full md:basis-1/2 wide:basis-1/4 p-2"
      />
    );
  });
  return <div className="flex flex-row flex-wrap">{rows}</div>;
}
