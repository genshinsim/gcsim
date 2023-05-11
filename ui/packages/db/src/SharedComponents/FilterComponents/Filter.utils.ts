import { createContext } from "react";
import { charNames } from "../../PipelineExtract/CharacterNames";

export interface FilterState {
  charFilter: CharFilter;
  charIncludeCount: number;
  pageNumber: number;
}

export enum ItemFilterState {
  "none",
  "include",
  "exclude",
}

export const initialCharFilter = charNames.reduce((acc, charName) => {
  acc[charName] = { state: ItemFilterState.none, charName };
  return acc;
}, {} as CharFilter);

export const FilterContext = createContext<FilterState>({
  charFilter: initialCharFilter,
  charIncludeCount: 0,
  pageNumber: 1,
});

// setName: number of pieces
// e.g. { "gladiatorsfinale": 2, "thundersoother": 4 }
export type ArtifactSetFilter = Record<string, number>;

export interface CharFilter {
  //character name
  [key: string]: CharFilterState;
}

export type CharFilterState =
  | {
      charName: string;
      state: ItemFilterState.include;
      weapon?: string;
      sets?: ArtifactSetFilter;
    }
  | {
      state: ItemFilterState.none;
      charName: string;
    }
  | {
      state: ItemFilterState.exclude;
      charName: string;
    };

export const FilterDispatchContext = createContext<
  React.Dispatch<FilterReducerAction>
>(null as unknown as React.Dispatch<FilterReducerAction>);

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
  filter: FilterState,
  action: FilterReducerAction
): FilterState {
  switch (action.type) {
    case "handleChar": {
      let newFilterState: ItemFilterState;
      let newCharIncludeCount = filter.charIncludeCount ?? 0;
      switch (filter.charFilter[action.char].state) {
        case ItemFilterState.none:
          //if more than 4 characters are included, do not include the new character
          if (filter.charIncludeCount >= 4)
            newFilterState = ItemFilterState.exclude;
          else {
            newFilterState = ItemFilterState.include;
            newCharIncludeCount++;
          }

          break;
        case ItemFilterState.include:
          newFilterState = ItemFilterState.exclude;
          newCharIncludeCount--;
          break;
        case ItemFilterState.exclude:
          newFilterState = ItemFilterState.none;
          break;
      }

      return {
        ...filter,
        charFilter: {
          ...filter.charFilter,
          [action.char]: {
            state: newFilterState,
            charName: action.char,
          },
        },
        charIncludeCount: newCharIncludeCount,
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
