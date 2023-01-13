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
    console.log(charFilter);
  }, [charFilter]);
  const determineButtonColor = (charName: string) => {
    const baseCss = "bp4-button ";
    switch (charFilter[charName]) {
      case FilterValue.none:
        return baseCss;
      case FilterValue.include:
        return baseCss + "bp4-intent-success";
      case FilterValue.exclude:
        return baseCss + "bp4-intent-danger";
      default:
        return baseCss;
    }
  };
  return (
    <div className="w-1/6 bg-slate-800 p-4">
      <div className="text-2xl pb-2 ">{t("Filter")}</div>
      <div
        className="flex flex-row  bp4-button"
        onClick={() => setCharIsOpen(!charIsOpen)}
      >
        <div className="grow ">{t("Character")}</div>
        <div>{charIsOpen ? "-" : "+"}</div>
      </div>
      <Collapse isOpen={charIsOpen}>
        <div className="flex flex-wrap gap-4 mt-2">
          {charNames.map((charName) => (
            <button
              className={determineButtonColor(charName)}
              onClick={() => {
                const newFilter = { ...charFilter };
                newFilter[charName] =
                  newFilter[charName] === FilterValue.none
                    ? FilterValue.include
                    : newFilter[charName] === FilterValue.include
                    ? FilterValue.exclude
                    : FilterValue.none;
                setCharFilter(newFilter);
              }}
            >
              {charName}
            </button>
          ))}
        </div>
      </Collapse>
    </div>
  );
}
