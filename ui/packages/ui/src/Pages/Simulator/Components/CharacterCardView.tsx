import { Character } from "@gcsim/types";
import { Trans, useTranslation } from "react-i18next";
import { CharacterCard } from "../../../Components/Cards";
import { ConsolidateCharStats } from "./character";

type Props = {
  chars: Character[];
  handleEdit: (index: number) => () => void;
};

export function CharacterCardView(props: Props) {
  useTranslation();

  if (!props.chars) {
    return (
      <div>
        <Trans>components.no_characters</Trans>
      </div>
    );
  }

  const teamStats = ConsolidateCharStats(props.chars);

  const rows = props.chars.map((c) => {
    return (
      <CharacterCard
        key={c.name}
        char={c}
        stats={teamStats.stats[c.name]}
        snapshot={teamStats.snapshot[c.name]}
        statsRows={teamStats.maxRows}
        handleDelete={() => console.log("deleting " + c.name)}
        className="basis-full md:basis-1/2 wide:basis-1/4 pt-2 pr-2 pb-2"
      />
    );
  });
  return <div className="flex flex-row flex-wrap pl-2">{rows}</div>;
}
