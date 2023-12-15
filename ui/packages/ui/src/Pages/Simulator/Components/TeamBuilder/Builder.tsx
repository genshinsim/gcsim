import { Icon } from "@blueprintjs/core";
import { Character } from "@gcsim/types";
import React from "react";
import { CharacterCard } from "../../../../Components/Cards";
import { ConsolidateCharStats } from "../character";

type Props = {
  team: Character[];
  handleAdd: () => void;
  handleRemove: (index: number) => () => void;
};

export const Builder = (props: Props) => {
  const [showDetails, setShowDetails] = React.useState(false);
  const teamStats = ConsolidateCharStats(props.team);

  // console.log(team);
  // console.log(teamStats);
  const handleToggleDetail = () => {
    setShowDetails(!showDetails);
  };

  const cards: JSX.Element[] = props.team.map((c, index) => {
    return (
      <CharacterCard
        key={c.name}
        char={c}
        stats={teamStats.stats[c.name]}
        statsRows={teamStats.maxRows}
        viewerMode
        handleToggleDetail={handleToggleDetail}
        showDetails={showDetails}
        handleDelete={props.handleRemove(index)}
        className="basis-full sm:basis-1/2 hd:basis-1/4 pt-2 pr-2 pb-2"
      />
    );
  });

  //add an extra card for adding new
  const blankCard = (
    <div
      className="basis-full sm:basis-1/2 hd:basis-1/4 pr-2 pb-2 pt-2"
      key="_blank"
    >
      <div
        className="bg-gray-600 shadow rounded-md hover:bg-gray-500 flex items-center justify-center min-h-[226px] h-full"
        onClick={props.handleAdd}
      >
        <Icon icon="plus" size={30} color="white" />
      </div>
    </div>
  );

  if (cards.length < 4) {
    cards.push(blankCard);
  }

  //TODO: add a button to toggle showing final stats
  return <div className="flex flex-row flex-wrap pl-2">{cards}</div>;
};
