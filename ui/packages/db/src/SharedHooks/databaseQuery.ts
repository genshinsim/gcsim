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
  const limit = entriesPerPage;
  const skip = (pageNumber - 1) * entriesPerPage;

  // sort all characters into included and excluded from the filter
  const includedChars: string[] = [];
  const excludedChars: string[] = [];
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
    const and: unknown[] = [];
    const trav: { [key in string]: boolean } = {};
    includedChars.forEach((char) => {
      if (char.includes("aether") || char.includes("lumine")) {
        const ele = char.replace(/(aether|lumine)(.+)/, "$2");
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
    const and: unknown[] = [];
    const trav: { [key in string]: boolean } = {};
    excludedChars.forEach((char) => {
      if (char.includes("aether") || char.includes("lumine")) {
        const ele = char.replace(/(aether|lumine)(.+)/, "$2");
        trav[ele] = true;
        return;
      }
      and.push({
        "summary.char_names": { $ne: char },
      });
    });
    Object.keys(trav).forEach((ele) => {
      and.push({
        "summary.char_names": { $ne: `aether${ele}` },
      });
      and.push({
        "summary.char_names": { $ne: `lumine${ele}` },
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

  // sort out tags
  const includedTags: string[] = [];
  for (const [tag, tagState] of Object.entries(filter.tagFilter)) {
    if (tagState.state === ItemFilterState.include) {
      includedTags.push(tag);
    }
  }

  if (includedTags.length > 0) {
    const tags: number[] = [];
    includedTags.forEach((tag) => {
      tags.push(parseInt(tag));
    });
    if (tags.length > 0) {
      query["accepted_tags"] = {
        $in: tags,
      };
    }
  }

  //force default sort by date for now
  return {
    query,
    limit,
    skip,
    sort: {
      create_date: -1,
    },
  };
}

export interface DbQuery {
  query: {
    $and?: unknown[];
    accepted_tags?: {
      $in?: number[];
    };
  };
  limit: number;
  sort?: {
    create_date: number;
  };
  skip?: number;
}
