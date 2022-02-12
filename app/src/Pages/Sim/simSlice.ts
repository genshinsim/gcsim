import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { Character, maxStatLength, Talent, Weapon } from "~/src/types";
import { characterKeyToICharacter } from "~src/Components/Character";
import { ascLvlMin, maxLvlToAsc } from "~src/util";

const testTeam: Character[] = [
  // {
  //   name: "bennett",
  //   element: "pyro",
  //   level: 70,
  //   max_level: 80,
  //   cons: 6,
  //   weapon: {
  //     name: "aquilafavonia",
  //     refine: 1,
  //     level: 90,
  //     max_level: 90,
  //   },
  //   talents: { attack: 9, skill: 9, burst: 10 },
  //   sets: { noblesseoblige: 4 },
  //   stats: [
  //     0, 0, 0, 5437, 0.146, 391, 0.105, 1.114, 46, 0.318, 0.98, 0, 0.466, 0, 0,
  //     0, 0, 0, 0, 0, 0, 0,
  //   ],
  //   snapshot: [
  //     0, 0, 0, 5437, 0.146, 391, 0.5549999999999999, 1.3806999995708467, 46,
  //     0.368, 1.48, 0, 0.466, 0, 0, 0, 0, 0, 0.41346000406980465, 0, 0, 0,
  //   ],
  // },
  // {
  //   name: "ganyu",
  //   element: "cryo",
  //   level: 90,
  //   max_level: 90,
  //   cons: 6,
  //   weapon: {
  //     name: "amosbow",
  //     refine: 5,
  //     level: 90,
  //     max_level: 90,
  //   },
  //   talents: { attack: 10, skill: 10, burst: 10 },
  //   sets: { wandererstroupe: 4 },
  //   stats: [
  //     0, 0.19, 0, 5527, 0, 365, 0.513, 0.175, 99, 0.459, 1.337, 0, 0, 0, 0.466,
  //     0, 0, 0, 0, 0, 0, 0,
  //   ],
  //   snapshot: [
  //     0, 0.19, 0, 5527, 0, 365, 1.2591519980381722, 0.175, 179, 0.509,
  //     2.22100000333786, 0, 0, 0, 0.466, 0, 0, 0, 0, 0, 0, 0,
  //   ],
  // },
  // {
  //   name: "kazuha",
  //   element: "anemo",
  //   level: 90,
  //   max_level: 90,
  //   cons: 0,
  //   weapon: {
  //     name: "freedomsworn",
  //     refine: 1,
  //     level: 90,
  //     max_level: 90,
  //   },
  //   talents: { attack: 9, skill: 9, burst: 9 },
  //   sets: { viridescentvenerer: 4 },
  //   stats: [
  //     0, 0.306, 70, 5557, 0, 393, 0, 0.777, 490, 0.10900000000000001, 0.621, 0,
  //     0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
  //   ],
  //   snapshot: [
  //     0, 0.306, 70, 5557, 0, 393, 0.25, 0.777, 803.6607945205687,
  //     0.15900000000000003, 1.121, 0, 0, 0, 0, 0, 0.15, 0, 0, 0, 0, 0.1,
  //   ],
  // },
  // {
  //   name: "xiangling",
  //   element: "pyro",
  //   level: 90,
  //   max_level: 90,
  //   cons: 6,
  //   weapon: {
  //     name: "favoniuslance",
  //     refine: 5,
  //     level: 90,
  //     max_level: 90,
  //   },
  //   talents: { attack: 2, skill: 9, burst: 10 },
  //   sets: { emblemofseveredfate: 4 },
  //   stats: [
  //     0, 0.241, 107, 5258, 0, 344, 0.099, 0.628, 102, 0.56, 0.823, 0, 0.466, 0,
  //     0, 0, 0, 0, 0, 0, 0, 0,
  //   ],
  //   snapshot: [
  //     0, 0.241, 107, 5258, 0, 344, 0.349, 1.1342681795149587, 198,
  //     0.6100000000000001, 1.323, 0, 0.466, 0, 0, 0, 0, 0, 0, 0, 0, 0,
  //   ],
  // },
];

export interface Sim {
  team: Character[];
  edit_index: number;
}

const initialState: Sim = {
  team: testTeam,
  edit_index: -1,
};

const defWep: { [key in string]: string } = {
  bow: "dullblade",
  catalyst: "dullblade",
  claymore: "dullblade",
  sword: "dullblade",
  polearm: "dullblade",
};

const newChar = (name: string): Character => {
  const c = characterKeyToICharacter[name];
  //default weapons
  return {
    name: name,
    level: 80,
    max_level: 90,
    element: c.element,
    cons: 0,
    weapon: {
      name: defWep[c.weapon_type],
      refine: 1,
      level: 1,
      max_level: 20,
    },
    talents: {
      attack: 6,
      skill: 6,
      burst: 6,
    },
    stats: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
    snapshot: [
      0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
    ],
    sets: {},
  };
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
    setCharacterNameAndEle: (
      state,
      action: PayloadAction<{ name: string; ele: string }>
    ) => {
      state.team[state.edit_index].name = action.payload.name;
      state.team[state.edit_index].element = action.payload.ele;
    },
    setCharacterLvl: (state, action: PayloadAction<{ val: number }>) => {
      state.team[state.edit_index].level = action.payload.val;
      return state;
    },
    setCharacterMaxLvl: (state, action: PayloadAction<{ val: number }>) => {
      let m = action.payload.val;
      let l = state.team[state.edit_index].level;
      let asc = maxLvlToAsc(m);
      if (l > m) {
        l = m;
      } else if (l < ascLvlMin(asc)) {
        l = ascLvlMin(asc);
      }

      state.team[state.edit_index].max_level = m;
      state.team[state.edit_index].level = l;
      return state;
    },
    setCharacterCon: (state, action: PayloadAction<{ val: number }>) => {
      state.team[state.edit_index].cons = action.payload.val;
      return state;
    },
    setCharacterSetBonus: (
      state,
      action: PayloadAction<{ set: string; val: number }>
    ) => {
      state.team[state.edit_index].sets[action.payload.set] =
        action.payload.val;
      return state;
    },
    addCharacterSet: (state, action: PayloadAction<{ set: string }>) => {
      state.team[state.edit_index].sets[action.payload.set] = 0;
      return state;
    },
    deleteCharacterSet: (state, action: PayloadAction<{ set: string }>) => {
      delete state.team[state.edit_index].sets[action.payload.set];
      return state;
    },
    setCharacterWeapon: (state, action: PayloadAction<{ val: Weapon }>) => {
      let w = action.payload.val;
      let asc = maxLvlToAsc(w.max_level);
      if (w.level > w.max_level) {
        w.level = w.max_level;
      } else if (w.level < ascLvlMin(asc)) {
        w.level = ascLvlMin(asc);
      }
      state.team[state.edit_index].weapon = w;
      return state;
    },
    setCharacterTalent: (state, action: PayloadAction<{ val: Talent }>) => {
      state.team[state.edit_index].talents = action.payload.val;
      return state;
    },
    setCharacterStats: (
      state,
      action: PayloadAction<{ index: number; val: number }>
    ) => {
      if (action.payload.index < 0 || action.payload.index > maxStatLength) {
        return state;
      }
      state.team[state.edit_index].stats[action.payload.index] =
        action.payload.val;
      return state;
    },
    addCharacter: (state, action: PayloadAction<{ name: string }>) => {
      if (state.team.length >= 4) return state;
      state.team.push(newChar(action.payload.name));
      return state;
    },
  },
});

export const simActions = simSlice.actions;

export type SimSlice = {
  [simSlice.name]: ReturnType<typeof simSlice["reducer"]>;
};
