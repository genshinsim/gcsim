import { Collapse, Drawer, DrawerSize, Position } from "@blueprintjs/core";
import { useContext, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { FaFilter, FaSearch } from "react-icons/fa";
import { charNames } from "../PipelineExtract/CharacterNames";
import useDebounce from "../SharedHooks/debounce";
import {
  FilterContext,
  FilterDispatchContext,
  ItemFilterState,
} from "./FilterComponents/Filter.utils";

export function Filter() {
  // https://github.com/i18next/next-i18next/issues/1795
  const { t: translation } = useTranslation();
  const t = (s: string) => translation<string>(s);

  const dispatch = useContext(FilterDispatchContext);
  const [isOpen, setIsOpen] = useState(false);

  // const filter = useContext(FilterContext);
  // const includedCharacterFilters: CharFilterState[] = Object.keys(
  //   filter.charFilter
  // )
  //   .filter((key) => filter.charFilter[key].state === ItemFilterState.include)
  //   .map((key) => filter.charFilter[key]);

  const [value, setValue] = useState<string>("");
  const debouncedValue = useDebounce<string>(value, 500);

  useEffect(() => {
    dispatch({ type: "setCustomFilter", customFilter: debouncedValue });
  }, [debouncedValue, dispatch]);

  return (
    <div>
      <button
        className="flex flex-row gap-2 bp4-button justify-center items-center p-3 bp4-intent-primary h-12 w-12"
        onClick={() => setIsOpen(!isOpen)}
      >
        <FaFilter size={24} className="opacity-80" />
      </button>

      <Drawer
        isOpen={isOpen}
        canEscapeKeyClose={true}
        canOutsideClickClose
        autoFocus
        isCloseButtonShown
        title={
          <div
            className="flex flex-row justify-between
          "
          >
            <div className="text-xl pb-1 ">{t("db.filter")}</div>
            <ClearFilterButton />
          </div>
        }
        onClose={() => setIsOpen(false)}
        position={Position.LEFT}
        size={DrawerSize.SMALL}
      >
        <div className="flex flex-col gap-2 overflow-y-auto overflow-x-hidden p-2">
          {/* <div className="flex flex-row gap-1">
            {includedCharacterFilters.map((charFilter) => (
              <FilterPortrait
                key={charFilter.charName}
                charName={charFilter.charName}
              />
            ))}
          </div> */}
          <input
            className="bp4-input bp4-icon bp4-icon-filter"
            placeholder={t("db.customFilter")}
            type="text"
            dir="auto"
            onChange={(e) => {
              setValue(e.target.value);
            }}
          />
          <CharacterFilter />
        </div>
      </Drawer>
    </div>
  );
}

function ClearFilterButton() {
  const { t: translation } = useTranslation();
  const t = (s: string) => translation<string>(s);
  const dispatch = useContext(FilterDispatchContext);
  return (
    <button
      className="bp4-button bp4-intent-danger bp4-small  "
      onClick={() => dispatch({ type: "clearFilter" })}
    >
      {t("db.clear")}
    </button>
  );
}

function CharacterFilter() {
  const [charIsOpen, setCharIsOpen] = useState(false);
  const { t: translation } = useTranslation();
  const t = (s: string) => translation<string>(s);
  const sortedCharNames = charNames.sort((a, b) => {
    if (t(a) < t(b)) {
      return -1;
    }
    if (t(a) > t(b)) {
      return 1;
    }
    return 0;
  });
  const [charSearch, setCharSearch] = useState<string>("");

  const translateCharName = (charName: string) =>
    t("game:character_names." + charName);

  return (
    <div className="w-full  overflow-x-hidden no-scrollbar">
      <button
        className=" bp4-button bp4-intent-primary w-full flex-row flex justify-between items-center "
        onClick={() => setCharIsOpen(!charIsOpen)}
      >
        <div className=" grow">{t("db.characters")}</div>

        <div className="">{charIsOpen ? "-" : "+"}</div>
      </button>
      <Collapse isOpen={charIsOpen}>
        <div className="flex flex-col mt-2 bg-gray-800 p-1">
          <label
            htmlFor="email"
            className="relative text-gray-400 focus-within:text-gray-600 flex flex-row"
          >
            <FaSearch className="pointer-events-none w-4 h-4 absolute top-2 transform   right-2 " />

            <input
              className="bp4-input bp4-icon bp4-icon-filter grow"
              type="text"
              dir="auto"
              onChange={(e) => {
                setCharSearch(e.target.value);
              }}
            />
          </label>

          <div className="grid grid-cols-4 gap-1 mt-1 overflow-y-auto overflow-x-hidden">
            {sortedCharNames
              .filter((charName) =>
                translateCharName(charName)
                  .toLocaleLowerCase()
                  .includes(charSearch.toLocaleLowerCase())
              )
              .map((charName) => (
                <CharFilterButton key={charName} charName={charName} />
              ))}
          </div>
        </div>
      </Collapse>
    </div>
  );
}

function CharFilterButton({ charName }: { charName: string }) {
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
          className={"bp4-button bp4-intent-success block"}
          onClick={handleClick}
        >
          <CharFilterButtonChild charName={charName} />
        </button>
      );
    case ItemFilterState.exclude:
      return (
        <button
          className={"bp4-button bp4-intent-danger block"}
          onClick={handleClick}
        >
          <CharFilterButtonChild charName={charName} />
        </button>
      );
    case ItemFilterState.none:
    default:
      return (
        <button className={"bp4-button block "} onClick={handleClick}>
          <CharFilterButtonChild charName={charName} />
        </button>
      );
  }
}

function CharFilterButtonChild({ charName }: { charName: string }) {
  const { t: translation } = useTranslation();
  const t = (s: string) => translation<string>(s);
  const displayCharName = t("game:character_names." + charName);

  return (
    <div className="flex flex-col truncate gap-1">
      <img
        alt={displayCharName}
        src={`/api/assets/avatar/${charName}.png`}
        className="truncate h-16 object-contain"
      />
      <div className="text-center">{displayCharName}</div>
    </div>
  );
}
