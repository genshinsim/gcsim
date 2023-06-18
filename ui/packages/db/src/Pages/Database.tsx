import { MenuItem, Spinner } from "@blueprintjs/core";
import { db } from "@gcsim/types";
import axios from "axios";
import { useContext, useEffect, useReducer, useState } from "react";
import { useTranslation } from "react-i18next";
import { newMockData } from "SharedComponents/mockData";
import { PaginationButtons } from "SharedComponents/Pagination";
import { Filter } from "../SharedComponents/Filter";
import {
  FilterContext,
  FilterDispatchContext,
  filterReducer,
  FilterState,
  initialFilter,
  ItemFilterState,
} from "../SharedComponents/FilterComponents/Filter.utils";
import { ListView } from "../SharedComponents/ListView";
import { MultiSelect2 } from "@blueprintjs/select";
import { charNames } from "PipelineExtract/CharacterNames";

export function Database() {
  const [filter, dispatch] = useReducer(filterReducer, initialFilter);

  const [data, setData] = useState<db.IEntry[]>([]);
  const { t } = useTranslation();
  const querydb = (query: DbQuery) => {
    axios(`/api/db?q=${encodeURIComponent(JSON.stringify(query))}`)
      .then((resp: { data: db.IEntries }) => {
        if (resp.data && resp.data.data) {
          setData(resp.data.data);
        } else {
          console.log("no data, using mockdata");
          setData(newMockData);
        }
      })
      .catch((err) => {
        console.log("error: ", err);
      });
  };

  useEffect(() => {
    const query = craftQuery(filter);
    querydb(query);
  }, [filter]);

  return (
    <FilterContext.Provider value={filter}>
      <FilterDispatchContext.Provider value={dispatch}>
        <div className="flex flex-col  gap-4 m-8 my-4 items-center">
          <div className="flex flex-row justify-between items-center w-full max-w-7xl ">
            <Filter />
            <CharacterQuickSelect />

            <div className="text-base  md:text-2xl">{`${t("db.showing")} ${
              data?.length ?? 0
            } ${t("db.simulations")} `}</div>
            {/* <Sorter /> */}
          </div>
          {data ? <ListView data={data} /> : <Spinner />}
          <PaginationButtons />
        </div>
      </FilterDispatchContext.Provider>
    </FilterContext.Provider>
  );
}

function craftQuery({
  charFilter,
  pageNumber,
  entriesPerPage,
  customFilter,
}: FilterState): DbQuery {
  const query: DbQuery["query"] = {};
  // sort all characters into included and excluded from the filter
  const includedChars: string[] = [];
  const excludedChars: string[] = [];
  const limit = entriesPerPage;
  const skip = (pageNumber - 1) * entriesPerPage;
  for (const [charName, charState] of Object.entries(charFilter)) {
    if (charState.state === ItemFilterState.include) {
      includedChars.push(charName);
    } else if (charState.state === ItemFilterState.exclude) {
      excludedChars.push(charName);
    }
  }

  if (customFilter) {
    let parsedFilter;
    try {
      parsedFilter = JSON.parse(`{${customFilter}}`);
    } catch (e) {
      console.log("invalid custom filter", e, customFilter);
    }

    return {
      query: parsedFilter,
      limit,
      skip,
    };
  }

  if (includedChars.length > 0) {
    query["summary.char_names"] = {};
    query["summary.char_names"]["$all"] = includedChars;
  }
  if (excludedChars.length > 0) {
    query["summary.char_names"] = query["summary.char_names"] ?? {};
    query["summary.char_names"]["$nin"] = excludedChars;
  }
  return {
    query,
    limit,
    skip,
  };
}

interface DbQuery {
  query: {
    "summary.char_names"?: {
      $all?: string[];
      $nin?: string[];
    };
  };
  limit?: number;
  sort?: unknown;
  skip?: number;
}

function CharacterQuickSelect() {
  //dispatch
  const dispatch = useContext(FilterDispatchContext);
  const filter = useContext(FilterContext);
  const { t } = useTranslation();

  const includedChars = Object.entries(filter.charFilter)
    .map(([charName, charState]) => {
      if (charState.state === ItemFilterState.include) {
        return charName;
      }
    })
    .filter((charName) => charName) as string[];

  const translateCharName = (charName: string) =>
    t("game:character_names." + charName);
  return (
    <div className="grow max-w-xl">
      <MultiSelect2
        items={charNames}
        itemRenderer={(charName, itemProps) => {
          return (
            <MenuItem
              key={charName}
              text={translateCharName(charName)}
              icon={
                <img
                  src={`/api/assets/avatar/${charName}.png`}
                  className="w-6 h-6"
                />
              }
              onClick={() => {
                dispatch({
                  type: "includeChar",
                  char: charName,
                });
              }}
              active={itemProps.modifiers.active}
            />
          );
        }}
        tagRenderer={(charName) => (
          <div className="flex flex-row gap-1" key={charName}>
            <img
              className="w-4 h-4"
              src={`/api/assets/avatar/${charName}.png`}
            />
            {translateCharName(charName)}
          </div>
        )}
        onItemSelect={(charName) => {
          if (!charName) {
            return;
          }
          dispatch({
            type: "handleChar",
            char: charName,
          });
        }}
        itemListPredicate={(query, items) => {
          return items.filter((item) => {
            return translateCharName(item)
              .toLocaleLowerCase()
              .includes(query.toLocaleLowerCase());
          });
        }}
        selectedItems={includedChars}
        onClear={() => {
          dispatch({
            type: "clearFilter",
          });
        }}
        onRemove={(charName) => {
          dispatch({
            type: "includeChar",
            char: charName,
          });
        }}
        resetOnSelect
        openOnKeyDown
        tagInputProps={{
          tagProps: {
            minimal: true,
          },
          onRemove: (value) => {
            if (!value) return;
            if (!value["key"]) return;
            dispatch({
              type: "removeChar",
              char: value["key"],
            });
          },
        }}
      ></MultiSelect2>
    </div>
  );
}
