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

  const cards = data.character_details.map((c) => {
    return (
      <CharacterCard
          key={c.name}
          char={c}
          showDetails={false}
          stats={[]}
          statsRows={0}
          className="basis-0 flex-auto min-w-[250px]" />
    );
  });

  return (
    <div className="col-span-full flex flex-row gap-2 justify-center flex-wrap">
      {cards}
    </div>
  );
};
