import { Collapse, Drawer, DrawerSize, Position } from "@blueprintjs/core";
import { useState } from "react";
import { FaFilter } from "react-icons/fa";
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

  const [isOpen, setIsOpen] = useState(false);
  return (
    <div className="">
      <button
        className="flex flex-row gap-2 bp4-button justify-center items-center p-3 bp4-intent-primary "
        onClick={() => setIsOpen(!isOpen)}
      >
        <FaFilter size={24} className="opacity-80" />
        {/* <div className="text-xl pb-1 ">{t("Filter")}</div> */}
      </button>
      <Drawer
        isOpen={isOpen}
        canEscapeKeyClose={true}
        canOutsideClickClose
        autoFocus
        isCloseButtonShown
        title={t("Filter")}
        onClose={() => setIsOpen(false)}
        position={Position.LEFT}
        size={DrawerSize.SMALL}
      >
        <CharacterFilter
          charFilter={charFilter}
          setCharFilter={setCharFilter}
        />
      </Drawer>
    </div>
  );
}

function CharacterFilter({
  charFilter,
  setCharFilter,
}: {
  charFilter: Record<string, FilterState>;
  setCharFilter: (newFilter: Record<string, FilterState>) => void;
}) {
  const [charIsOpen, setCharIsOpen] = useState(false);
  const t = useTranslation;

  return (
    <div className="w-full">
      <button
        className=" bp4-button bp4-intent-primary pl-5 pr-3 w-full "
        onClick={() => setCharIsOpen(!charIsOpen)}
      >
        <div className=" grow">{t("Character")}</div>
        <div className="">{charIsOpen ? "-" : "+"}</div>
      </button>
      <Collapse isOpen={charIsOpen}>
        <div className="grid grid-cols-3 gap-1 mt-1 ">
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
