import { Collapse } from "@blueprintjs/core";
import { useEffect, useState } from "react";

const charNames = [
  //genshin characters
  "albedo",
  "amber",
  "barbara",
  "beidou",
  "bennett",
  "chongyun",
  "diluc",
  "diona",
  "fischl",
  "ganyu",
  "hutao",
  "jean",
  "kaeya",
  "keqing",
  "klee",
  "lisa",
  "mona",
  "ningguang",
  "noelle",
  "qiqi",
  "razor",
  "rosaria",
  "sucrose",
  "tartaglia",
  "traveler",
  "venti",
  "xiangling",
  "xingqiu",
  "xinyan",
  "xiao",
  "yanfei",
  "zhongli",
  "nahida",
  "yelan",
  "raiden",
];

const useTranslation = (key: string) => key;
enum FilterValue {
  "none",
  "include",
  "exclude",
}
export function Filter() {
  const t = useTranslation;
  const [charIsOpen, setCharIsOpen] = useState(false);
  const [charFilter, setCharFilter] = useState<Record<string, FilterValue>>(
    //use charNames to create an object with all characters as keys and empty strings as values
    charNames.reduce((acc, charName) => {
      acc[charName] = FilterValue.none;
      return acc;
    }, {} as Record<string, FilterValue>)
  );

  useEffect(() => {
    // console.log(charFilter);
  }, [charFilter]);

  return (
    <div className="w-96 bg-slate-800 p-4">
      <div className="text-2xl pb-2 ">{t("Filter")}</div>
      <div
        className="flex flex-row  bp4-button"
        onClick={() => setCharIsOpen(!charIsOpen)}
      >
        <div className="grow ">{t("Character")}</div>
        <div>{charIsOpen ? "-" : "+"}</div>
      </div>
      <Collapse isOpen={charIsOpen}>
        <div className="grid grid-cols-3 gap-2 mt-2">
          {charNames.map((charName) => (
            <CharFilterButton
              charName={charName}
              charFilter={charFilter}
              setCharFilter={setCharFilter}
            />
          ))}
        </div>
      </Collapse>
    </div>
  );
}

function CharFilterButton({
  charName,
  charFilter,
  setCharFilter,
}: {
  charName: string;
  charFilter: Record<string, FilterValue>;
  setCharFilter: (newFilter: Record<string, FilterValue>) => void;
}) {
  const t = useTranslation;
  const displayCharName = t(charName);

  const handleClick = () => {
    const newFilter = { ...charFilter };
    newFilter[charName] =
      newFilter[charName] === FilterValue.none
        ? FilterValue.include
        : newFilter[charName] === FilterValue.include
        ? FilterValue.exclude
        : FilterValue.none;
    setCharFilter(newFilter);
  };

  switch (charFilter[charName]) {
    case FilterValue.include:
      return (
        <button
          className={"bp4-button bp4-intent-success"}
          onClick={handleClick}
        >
          {`+ ` + displayCharName}
        </button>
      );
    case FilterValue.exclude:
      return (
        <button
          className={"bp4-button bp4-intent-danger"}
          onClick={handleClick}
        >
          {`- ` + displayCharName}
        </button>
      );
    case FilterValue.none:
    default:
      return (
        <button className={"bp4-button"} onClick={handleClick}>
          {displayCharName}
        </button>
      );
  }
}
