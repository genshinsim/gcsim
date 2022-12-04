import { SimResults } from "@gcsim/types";
import classNames from "classnames";
import { CharacterCard } from "../../../../Components/Cards";

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

const CharacterCards = ({ data }: Props) => {
  const cardClass = characterCardsClassNames(data?.character_details?.length ?? 4);

  if (data?.character_details == null) {
    return (
      <>
        <FakeCard className={cardClass} />
        <FakeCard className={cardClass} />
        <FakeCard className={cardClass} />
        <FakeCard className={cardClass} />
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