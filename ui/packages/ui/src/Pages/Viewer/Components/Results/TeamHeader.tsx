import { SimResults } from "@gcsim/types";
import { CharacterCard } from "../../../../Components/Cards";

const cardClass = "basis-0 flex-auto min-w-[250px]";

type Props = {
  data: SimResults | null;
};

export default ({ data }: Props) => {
  return (
    <div className="col-span-full flex flex-row gap-2 justify-center flex-wrap">
      <CharacterCards data={data} />
    </div>
  );
};


const CharacterCards = ({ data }: Props) => {
  if (data?.character_details == null) {
    return (
      <>
        <FakeCard />
        <FakeCard />
        <FakeCard />
        <FakeCard />
      </>
    );
  }

  return (
    <>
    {data.character_details.map((c) => (
      <CharacterCard
        key={c.name}
        char={c}
        showDetails={false}
        stats={[]}
        statsRows={0}
        className={cardClass} />
    ))}
    </>
  );
};

const FakeCard = ({}) => (
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
      className={cardClass} />
);