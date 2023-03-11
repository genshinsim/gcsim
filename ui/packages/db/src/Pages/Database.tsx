import { Spinner } from "@blueprintjs/core";
import { model } from "@gcsim/types";
import { useEffect, useReducer, useState } from "react";
import { Filter } from "../SharedComponents/Filter";
import {
  CharFilter,
  FilterContext,
  FilterDispatchContext,
  filterReducer,
  FilterState,
  initialCharFilter,
} from "../SharedComponents/FilterComponents/Filter.utils";
import { ListView } from "../SharedComponents/ListView";
import { mockData } from "../SharedComponents/mockData";
import Sorter from "../SharedComponents/Sorter";

export function Database() {
  const [filter, dispatch] = useReducer(filterReducer, {
    charFilter: initialCharFilter,
  });

  const [data, setData] = useState<model.IDBEntries["data"]>([]);

  useEffect(() => {
    //https://simimpact.app/api
    const url = `api/db?q=${encodeURIComponent(
      JSON.stringify(filter.charFilter)
    )}`;
    fetch(url)
      .then((res) => res.json())
      .then((data) => {
        // setData(data.data);
        setData(mockData);
      })
      .catch((e) => {
        console.log(e);
      });
  }, [filter.charFilter]);

  return (
    <FilterContext.Provider value={filter}>
      <FilterDispatchContext.Provider value={dispatch}>
        <div className="flex flex-col  gap-4 m-8 my-4 ">
          <div className="flex flex-row justify-between items-center">
            <Filter />
            <div className="text-base  md:text-2xl">{`Showing ${
              data?.length ?? 0
            } Simulations `}</div>
            <Sorter />
          </div>
          {data ? <ListView data={data} /> : <Spinner />}
        </div>
      </FilterDispatchContext.Provider>
    </FilterContext.Provider>
  );
}

function craftQuery({ charFilter }: CharFilter): unknown {
  const query: Record<string, unknown> = {};
  // sort all characters into included and excluded from the filter
  const includedChars: string[] = [];
  const excludedChars: string[] = [];
  for (const [charName, filterState] of Object.entries(charFilter)) {
    if (filterState === FilterState.include) {
      includedChars.push(charName);
    } else if (filterState === FilterState.exclude) {
      excludedChars.push(charName);
    }
  }
  if (includedChars.length > 0 || excludedChars.length > 0) {
    query.char_names = { $all: includedChars, $nin: excludedChars };
  }
  return query;
}
