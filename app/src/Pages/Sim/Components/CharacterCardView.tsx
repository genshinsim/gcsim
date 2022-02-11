import {
  CharDetail,
  ConsolidateCharStats,
  CharacterCard,
} from "~src/Components/Character";

type Props = {
  chars: CharDetail[];
  handleEdit: (index: number) => () => void;
};

export function CharacterCardView(props: Props) {
  const teamStats = ConsolidateCharStats(props.chars);

  const rows = props.chars.map((c, index) => {
    return (
      <CharacterCard
        key={c.name}
        char={c}
        stats={teamStats.stats[c.name]}
        statsRows={teamStats.maxRows}
        showDelete
        showEdit
        handleDelete={() => console.log("deleting " + c.name)}
        toggleEdit={props.handleEdit(index)}
        className="basis-full md:basis-1/2 wide:basis-1/4 pt-2 pr-2 pb-2"
      />
    );
  });
  return <div className="flex flex-row flex-wrap pl-2">{rows}</div>;
}
