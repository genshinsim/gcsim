import { SimResults } from "@gcsim/types";
import { CharacterCard } from "../../../../Components/Cards";

type Props = {
  data: SimResults | null;
};

export default ({ data }: Props) => {
  if (data === null) {
    return null;
  }

  if (!data.character_details) {
    return null;
  }

  const cards: JSX.Element[] = data.character_details.map((c, index) => {
    return (
      <CharacterCard
        key={c.name}
        char={c}
        showDetails={false}
        stats={[]}
        statsRows={0}
        className="basis-full sm:basis-1/2 hd:basis-1/4 pt-2 pr-2 pb-2"
      />
    );
  });

  return (
    <div className="m-w-full flex flex-row flex-wrap -mr-2 justify-center">
      {cards}
    </div>
  );
};
