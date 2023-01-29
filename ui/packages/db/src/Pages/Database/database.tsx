import { Filter, FilterValue } from "./Filter";
import { ListView, ListViewProps } from "./ListView";
import { useEffect, useState } from "react";
import { charNames } from "../../PipelineExtract/CharacterNames.";

export function Database() {
  const [charFilter, setCharFilter] = useState<Record<string, FilterValue>>(
    //use charNames to create an object with all characters as keys and empty strings as values
    charNames.reduce((acc, charName) => {
      acc[charName] = FilterValue.none;
      return acc;
    }, {} as Record<string, FilterValue>)
  );
  const [key, setKey] = useState(0);

  useEffect(() => {
    setKey((prev) => prev + 1);
  }, [charFilter]);

  return (
    <div className="flex flex-row gap-4">
      <Filter charFilter={charFilter} setCharFilter={setCharFilter} />
      <ListView query={craftQuery(charFilter)} key={key} />
    </div>
  );
}

function craftQuery(
  charFilter: Record<string, FilterValue>
): ListViewProps["query"] {
  const query: Record<string, any> = {};
  const charNamesArray = Object.entries(charFilter)
    .filter(([charName, value]) => {
      return value === FilterValue.include;
    })
    .map(([charName, value]) => charName);
  if (charNamesArray.length > 0) {
    query.char_names = { $all: charNamesArray };
  }
  return query;
}
