import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { Character } from "~/src/types";

const testTeam: Character[] = [
  {
    name: "bennett",
    element: "pyro",
    level: 70,
    max_level: 80,
    cons: 6,
    weapon: {
      name: "aquilafavonia",
      refine: 1,
      level: 90,
      max_level: 90,
    },
    talents: { attack: 9, skill: 9, burst: 10 },
    sets: { noblesseoblige: 4 },
    stats: [
      0, 0, 0, 5437, 0.146, 391, 0.105, 1.114, 46, 0.318, 0.98, 0, 0.466, 0, 0,
      0, 0, 0, 0, 0, 0, 0,
    ],
    snapshot: [
      0, 0, 0, 5437, 0.146, 391, 0.5549999999999999, 1.3806999995708467, 46,
      0.368, 1.48, 0, 0.466, 0, 0, 0, 0, 0, 0.41346000406980465, 0, 0, 0,
    ],
  },
  {
    name: "ganyu",
    element: "cryo",
    level: 90,
    max_level: 90,
    cons: 6,
    weapon: {
      name: "amosbow",
      refine: 5,
      level: 90,
      max_level: 90,
    },
    talents: { attack: 10, skill: 10, burst: 10 },
    sets: { wandererstroupe: 4 },
    stats: [
      0, 0.19, 0, 5527, 0, 365, 0.513, 0.175, 99, 0.459, 1.337, 0, 0, 0, 0.466,
      0, 0, 0, 0, 0, 0, 0,
    ],
    snapshot: [
      0, 0.19, 0, 5527, 0, 365, 1.2591519980381722, 0.175, 179, 0.509,
      2.22100000333786, 0, 0, 0, 0.466, 0, 0, 0, 0, 0, 0, 0,
    ],
  },
  {
    name: "kazuha",
    element: "anemo",
    level: 90,
    max_level: 90,
    cons: 0,
    weapon: {
      name: "freedomsworn",
      refine: 1,
      level: 90,
      max_level: 90,
    },
    talents: { attack: 9, skill: 9, burst: 9 },
    sets: { viridescentvenerer: 4 },
    stats: [
      0, 0.306, 70, 5557, 0, 393, 0, 0.777, 490, 0.10900000000000001, 0.621, 0,
      0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
    ],
    snapshot: [
      0, 0.306, 70, 5557, 0, 393, 0.25, 0.777, 803.6607945205687,
      0.15900000000000003, 1.121, 0, 0, 0, 0, 0, 0.15, 0, 0, 0, 0, 0.1,
    ],
  },
  {
    name: "xiangling",
    element: "pyro",
    level: 90,
    max_level: 90,
    cons: 6,
    weapon: {
      name: "favoniuslance",
      refine: 5,
      level: 90,
      max_level: 90,
    },
    talents: { attack: 2, skill: 9, burst: 10 },
    sets: { emblemofseveredfate: 4 },
    stats: [
      0, 0.241, 107, 5258, 0, 344, 0.099, 0.628, 102, 0.56, 0.823, 0, 0.466, 0,
      0, 0, 0, 0, 0, 0, 0, 0,
    ],
    snapshot: [
      0, 0.241, 107, 5258, 0, 344, 0.349, 1.1342681795149587, 198,
      0.6100000000000001, 1.323, 0, 0.466, 0, 0, 0, 0, 0, 0, 0, 0, 0,
    ],
  },
];

export interface Sim {
  team: Character[];
  edit_index: number;
}

const initialState: Sim = {
  team: testTeam,
  edit_index: -1,
};

export const simSlice = createSlice({
  name: "sim",
  initialState: initialState,
  reducers: {
    deleteCharacter: (state, action: PayloadAction<{ index: number }>) => {
      if (
        action.payload.index < 0 ||
        action.payload.index >= state.team.length
      ) {
        return state;
      }
      state.team.splice(action.payload.index, 1);
    },
    editCharacter: (state, action: PayloadAction<{ index: number }>) => {
      //if index < 0 then it means it's closed; but it shouldn't be bigger than max length
      if (action.payload.index >= state.team.length) {
        return state;
      }
      state.edit_index = action.payload.index;
    },
    setCharacter: (
      state,
      action: PayloadAction<{ char: Character; index: number }>
    ) => {
      if (
        action.payload.index < 0 ||
        action.payload.index >= state.team.length
      ) {
        return state;
      }
      state.team[action.payload.index] = action.payload.char;
      return state;
    },
  },
});

export const simActions = simSlice.actions;

export type SimSlice = {
  [simSlice.name]: ReturnType<typeof simSlice["reducer"]>;
};
