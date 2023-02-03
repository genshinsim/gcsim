import { Collapse } from "@blueprintjs/core";
import { useState } from "react";
import { charNames } from "../../../PipelineExtract/CharacterNames.";

const useTranslation = (key: string) => key;
export enum FilterState {
  "none",
  "include",
  "exclude",
}
export function Filter({
  charFilter,
  setCharFilter,
}: {
  charFilter: Record<string, FilterState>;
  setCharFilter: (newFilter: Record<string, FilterState>) => void;
}) {
  const t = useTranslation;
  const [charIsOpen, setCharIsOpen] = useState(false);

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
        <div className="grid grid-cols-3 gap-1 mt-1">
          {charNames.map((charName) => (
            <CharFilterButton
              key={charName}
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
  charFilter: Record<string, FilterState>;
  setCharFilter: (newFilter: Record<string, FilterState>) => void;
}) {
  const t = useTranslation;
  const displayCharName = t(charName);

  const handleClick = () => {
    const newFilter = { ...charFilter };
    newFilter[charName] =
      newFilter[charName] === FilterState.none
        ? FilterState.include
        : newFilter[charName] === FilterState.include
        ? FilterState.exclude
        : FilterState.none;
    setCharFilter(newFilter);
  };

  switch (charFilter[charName]) {
    case FilterState.include:
      return (
        <button
          className={"bp4-button bp4-intent-success"}
          onClick={handleClick}
        >
          {`+ ` + displayCharName}
        </button>
      );
    case FilterState.exclude:
      return (
        <button
          className={"bp4-button bp4-intent-danger"}
          onClick={handleClick}
        >
          {`- ` + displayCharName}
        </button>
      );
    case FilterState.none:
    default:
      return (
        <button className={"bp4-button"} onClick={handleClick}>
          {displayCharName}
        </button>
      );
  }
}
