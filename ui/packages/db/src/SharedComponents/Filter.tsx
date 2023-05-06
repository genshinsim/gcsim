import { Collapse, Drawer, DrawerSize, Position } from "@blueprintjs/core";
import { useContext, useState } from "react";
import { FaFilter } from "react-icons/fa";
import { charNames } from "../PipelineExtract/CharacterNames";
import {
  FilterContext,
  FilterDispatchContext,
  ItemFilterState,
} from "./FilterComponents/Filter.utils";

const useTranslation = (key: string) => key;

export function Filter() {
  const t = useTranslation;

  const [isOpen, setIsOpen] = useState(false);

  // const filter = useContext(FilterContext);
  // const includedCharacterFilters: CharFilterState[] = Object.keys(
  //   filter.charFilter
  // )
  //   .filter((key) => filter.charFilter[key].state === ItemFilterState.include)
  //   .map((key) => filter.charFilter[key]);

  return (
    <div>
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
        <div className="flex flex-col gap-2 overflow-y-auto overflow-x-hidden">
          {/* <div className="flex flex-row gap-1">
            {includedCharacterFilters.map((charFilter) => (
              <FilterPortrait
                key={charFilter.charName}
                charName={charFilter.charName}
              />
            ))}
          </div> */}
          <CharacterFilter />
        </div>
      </Drawer>
    </div>
  );
}

// function FilterDrawer(charFilter: Record<string, FilterState>) {
//   return (
//     <div className="w-full overflow-y-auto overflow-x-hidden no-scrollbar"></div>
//   );
// }

function CharacterFilter() {
  const [charIsOpen, setCharIsOpen] = useState(false);
  const t = useTranslation;

  return (
    <div className="w-full  overflow-x-hidden no-scrollbar">
      <button
        className=" bp4-button bp4-intent-primary pl-5 pr-3 w-full "
        onClick={() => setCharIsOpen(!charIsOpen)}
      >
        <div className=" grow">{t("Characters")}</div>
        <div className="">{charIsOpen ? "-" : "+"}</div>
      </button>
      <Collapse isOpen={charIsOpen}>
        <div className="grid grid-cols-3 gap-1 mt-1 overflow-y-auto overflow-x-hidden">
          {charNames.map((charName) => (
            <FilterButton key={charName} charName={charName} />
          ))}
        </div>
      </Collapse>
    </div>
  );
}

function FilterButton({ charName }: { charName: string }) {
  const t = useTranslation;
  const displayCharName = t(charName);
  const filter = useContext(FilterContext);
  const dispatch = useContext(FilterDispatchContext);

  const handleClick = () => {
    dispatch({
      type: "handleChar",
      char: charName,
    });
  };

  switch (filter.charFilter[charName].state) {
    case ItemFilterState.include:
      return (
        <button
          className={"bp4-button bp4-intent-success"}
          onClick={handleClick}
        >
          {`+ ` + displayCharName}
        </button>
      );
    case ItemFilterState.exclude:
      return (
        <button
          className={"bp4-button bp4-intent-danger"}
          onClick={handleClick}
        >
          {`- ` + displayCharName}
        </button>
      );
    case ItemFilterState.none:
    default:
      return (
        <button className={"bp4-button"} onClick={handleClick}>
          {displayCharName}
        </button>
      );
  }
}
