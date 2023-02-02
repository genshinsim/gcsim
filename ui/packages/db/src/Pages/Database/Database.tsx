import { model } from "@gcsim/types";

import { Spinner } from "@blueprintjs/core";
import { useEffect, useState } from "react";
import { charNames } from "../../PipelineExtract/CharacterNames.";
import { mockData } from "./Components/mockData";
import { Filter, FilterState } from "./Filter";
import { ListViewProps } from "./ListView";

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
    const url = `https://simimpact.app/api/db?q=${encodeURIComponent(
      JSON.stringify(craftQuery(charFilter))
    )}`;
    fetch(url)
      .then((res) => res.json())
      .then((data) => {
        console.log(data);
        setData(mockData);
      })
      .catch((e) => {
        console.log(e);
      });
  }, [charFilter]);

  if (!data) {
    return (
      <div>
        <Spinner />
      </div>
    );
  }

  return (
    <div className="flex flex-row gap-4">
      <Filter charFilter={charFilter} setCharFilter={setCharFilter} />
      {/* <ListView query={craftQuery(charFilter)} key={key} /> */}
    </div>
  );
}

function craftQuery(
  charFilter: Record<string, FilterState>
): ListViewProps["query"] {
  const query: Record<string, any> = {};
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
