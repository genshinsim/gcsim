import { Character } from "@gcsim/types";
import classNames from "classnames";
import React, { memo } from "react";
import { CharacterCard } from "../../../../Components/Cards";
import { ConsolidateCharStats } from "../../../Simulator/Components/character";

type Props = {
  characters?: Character[];
};

const TeamHeader = ({ characters }: Props) => {
  return (
    <div className="col-span-full flex flex-row gap-2 justify-center flex-wrap">
      <CharacterCards characters={characters} />
    </div>
  );
};

export default memo(TeamHeader);

export function characterCardsClassNames(num: number): string {
  return classNames(
    "basis-0 flex-auto",
    {
      ["min-[300px]:min-w-[300px]"]: num % 2 == 1, // for special 3 char case
      ["min-[270px]:min-w-[270px]"]: num % 2 == 0,
      ["min-[825px]:min-w-[400px]"]: num % 2 == 0,
      ["min-[1200px]:min-w-[270px]"]: num % 2 == 0,
  });
}

const CharacterCards = ({ characters }: Props) => {
  const cardClass = characterCardsClassNames(characters?.length ?? 4);
  const [showDetails, setShowDetails] = React.useState(false);

  if (characters == null) {
    return (
      <>
        <FakeCard className={cardClass} />
        <FakeCard className={cardClass} />
        <FakeCard className={cardClass} />
        <FakeCard className={cardClass} />
      </>
    );
  }

  const handleToggleDetail = () => {
    setShowDetails(!showDetails)
  }

  const statBlock = ConsolidateCharStats(characters)

  console.log(characters)

  return (
    <>
    {characters.map((c) => (
      <CharacterCard
        key={c.name}
        char={c}
        showDetails={showDetails}
        handleToggleDetail={handleToggleDetail}
        viewerMode
        stats={statBlock.stats[c.name] ? statBlock.stats[c.name] : []}
        statsRows={statBlock.maxRows ? statBlock.maxRows : 0}
        className={cardClass} />
    ))}
    </>
  );
};

export const FakeCard = ({ className }: { className: string }) => (
  <CharacterCard
      key="fake"
      char={{
        name: "fake name",
        level: 90,
        max_level: 90,
        element: "none",
        cons: 6,
        weapon: {
          name: "fake weapon",
          refine: 6,
          level: 90,
          max_level: 90
        },
        talents: {
          attack: 9,
          skill: 9,
          burst: 9
        },
        stats: [],
        snapshot: [],
        sets: {}
      }}
      showDetails={false}
      stats={[]}
      statsRows={0}
      isSkeleton={true}
      className={className} />
);