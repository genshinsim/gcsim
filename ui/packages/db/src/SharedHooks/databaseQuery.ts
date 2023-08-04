import {
  FilterState,
  ItemFilterState,
} from "SharedComponents/FilterComponents/Filter.utils";

export function craftQuery(
  filter: FilterState,
  pageNumber: number,
  entriesPerPage: number
): DbQuery {
  const query: DbQuery["query"] = {};
  // sort all characters into included and excluded from the filter
  const includedChars: string[] = [];
  const excludedChars: string[] = [];
  const limit = entriesPerPage;
  const skip = (pageNumber - 1) * entriesPerPage;
  for (const [charName, charState] of Object.entries(filter.charFilter)) {
    if (charState.state === ItemFilterState.include) {
      includedChars.push(charName);
    } else if (charState.state === ItemFilterState.exclude) {
      excludedChars.push(charName);
    }
  }

  if (filter.customFilter) {
    let parsedFilter;
    try {
      parsedFilter = JSON.parse(`{${filter.customFilter}}`);
    } catch (e) {
      console.log("invalid custom filter", e, filter.customFilter);
    }

    return {
      query: parsedFilter,
      limit,
      skip,
    };
  }

  if (includedChars.length > 0) {
    let and: any[] = [];
    let trav: { [key in string]: boolean } = {};
    includedChars.forEach((char) => {
      if (char.includes("aether") || char.includes("lumine")) {
        let ele = char.replace(/(aether|lumine)(.+)/, "$2");
        trav[ele] = true;
        return;
      }
      and.push({
        "summary.char_names": char,
      });
    });
    Object.keys(trav).forEach((ele) => {
      and.push({
        $or: [
          { "summary.char_names": `aether${ele}` },
          { "summary.char_names": `lumine${ele}` },
        ],
      });
    });
    if (and.length > 0) {
      query["$and"] = and;
    }
  }

  if (excludedChars.length > 0) {
    let and: any[] = [];
    let trav: { [key in string]: boolean } = {};
    excludedChars.forEach((char) => {
      if (char.includes("aether") || char.includes("lumine")) {
        let ele = char.replace(/(aether|lumine)(.+)/, "$2");
        trav[ele] = true;
        return;
      }
      and.push({
          "summary.char_names": {$ne: char},
      });
    });
    Object.keys(trav).forEach((ele) => {
      and.push({
          "summary.char_names": {$ne: `aether${ele}`},
      });
      and.push({
          "summary.char_names": {$ne: `lumine${ele}`},
      });
    });
    if (and.length > 0) {
      if (query["$and"] === undefined) {
        query["$and"] = [];
      }
      and.forEach((e) => {
        query["$and"]?.push(e);
      });
    }
  }
  if (filter.tags.length > 0) {
    query["accepted_tags"] = {
      $in: filter.tags,
    };
  }
  return {
    query,
    limit,
    skip,
  };
}

export interface DbQuery {
  query: {
    $and?: any[];
    accepted_tags?: {
      $in?: number[];
    };
  };
  limit: number;
  sort?: unknown;
  skip?: number;
}
