import { Spinner } from "@blueprintjs/core";
import { model } from "@gcsim/types";
import { useEffect, useState } from "react";
import { charNames } from "../../PipelineExtract/CharacterNames.";
import { Filter, FilterState } from "./Components/Filter";
import { ListView } from "./Components/ListView";
import { mockData } from "./Components/mockData";
import Sorter from "./Components/Sorter";

export function Database() {
  const [charFilter, setCharFilter] = useState<Record<string, FilterState>>(
    //use charNames to create an object with all characters as keys and empty strings as values for default
    charNames.reduce((acc, charName) => {
      acc[charName] = FilterState.none;
      return acc;
    }, {} as Record<string, FilterState>)
  );
  const [data, setData] = useState<model.IDBEntries["data"]>([]);

  useEffect(() => {
    //https://simimpact.app/api
    const url = `api/db?q=${encodeURIComponent(
      JSON.stringify(craftQuery(charFilter))
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
  }, [charFilter]);

  return (
    <div className="flex flex-col  gap-4 m-8 my-4">
      <div className="flex flex-row justify-between items-center">
        <Filter charFilter={charFilter} setCharFilter={setCharFilter} />
        <div className="text-base  md:text-2xl">{`Showing ${
          data?.length ?? 0
        } Simulations `}</div>
        <Sorter />
      </div>
      {data ? <ListView data={data} /> : <Spinner />}
    </div>
  );
}

function craftQuery(charFilter: Record<string, FilterState>): unknown {
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
