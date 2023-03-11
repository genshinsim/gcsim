import { createContext } from "react";

export enum FilterState {
  "none",
  "include",
  "exclude",
}

export const FilterContext = createContext<{
  charFilter: Record<string, FilterState>;
}>({
  charFilter: {},
});
