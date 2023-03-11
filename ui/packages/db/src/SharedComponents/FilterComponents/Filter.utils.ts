import { createContext } from "react";
import { charNames } from "../../PipelineExtract/CharacterNames";

export enum FilterState {
  "none",
  "include",
  "exclude",
}

export interface CharFilter {
  //character name
  [key: string]: CharFilterState;
}

export type CharFilterState =
  | {
      state: FilterState.include;
      weapon: string;
      set: {
        //set name
        [key: string]: number;
      };
    }
  | {
      state: FilterState.none;
    }
  | {
      state: FilterState.exclude;
    };

export const FilterDispatchContext = createContext<
  React.Dispatch<FilterReducerAction>
>(null as unknown as React.Dispatch<FilterReducerAction>);

export const initialCharFilter = charNames.reduce((acc, char) => {
  acc[char] = { state: FilterState.none };
  return acc;
}, {} as CharFilter);

export const FilterContext = createContext<{
  charFilter: CharFilter;
}>({ charFilter: initialCharFilter });

export interface FilterReducerAction {
  type:
    | "handleChar"
    | "includeWeapon"
    | "nullWeapon"
    | "includeSet"
    | "nullSet";
  char: string;
  weapon?: string;
  set?: string;
}
export function filterReducer(
  filter: {
    charFilter: CharFilter;
  },
  action: FilterReducerAction
): { charFilter: CharFilter } {
  console.log(action);
  console.log("real", filter);
  switch (action.type) {
    case "handleChar": {
      let newFilterState;
      switch (filter.charFilter[action.char].state) {
        case FilterState.none:
          newFilterState = FilterState.include;
          break;
        case FilterState.include:
          newFilterState = FilterState.exclude;
          break;
        case FilterState.exclude:
          newFilterState = FilterState.none;
          break;
      }
      return {
        ...filter,
        charFilter: {
          ...filter.charFilter,
          [action.char]: {
            state: newFilterState,
          },
        },
      };
    }
    case "includeWeapon": {
      return {
        ...filter,
        [action.char]: {
          ...filter[action.char],
          weapon: action.weapon,
        },
      };
    }

    case "nullWeapon": {
      return {
        ...filter,
        [action.char]: {
          ...filter[action.char],
          weapon: "",
        },
      };
    }
    case "includeSet": {
      if (!filter[action.char].set) filter[action.char].set = {};
      if (filter[action.char].set[action.set as string] === 2) {
        return {
          ...filter,
          [action.char]: {
            ...filter[action.char],
            set: {
              ...filter[action.char].set,
              [action.set as string]: 2,
            },
          },
        };
      }
      return {
        ...filter,
        [action.char]: {
          ...filter[action.char],
          set: {
            ...filter[action.char].set,
            [action.set as string]: 2,
          },
        },
      };
    }
    case "nullSet": {
      const { [action.set as string]: _, ...newSet } = filter[action.char].set;
      return {
        ...filter,
        [action.char]: {
          ...filter[action.char],
          set: newSet,
        },
      };
    }

    default: {
      throw Error("Unknown action: " + action.type);
    }
  }
}
