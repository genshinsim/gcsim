import {
  FilterState,
  ItemFilterState,
} from "SharedComponents/FilterComponents/Filter.utils";

export function craftQuery(filter: FilterState, pageNumber: number, entriesPerPage: number): DbQuery {
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
    query["summary.char_names"] = {};
    query["summary.char_names"]["$all"] = includedChars;
  }
  if (excludedChars.length > 0) {
    query["summary.char_names"] = query["summary.char_names"] ?? {};
    query["summary.char_names"]["$nin"] = excludedChars;
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
    "summary.char_names"?: {
      $all?: string[];
      $nin?: string[];
    };
    accepted_tags?: {
      $in?: number[];
    };
  };
  limit: number;
  sort?: unknown;
  skip?: number;
}
