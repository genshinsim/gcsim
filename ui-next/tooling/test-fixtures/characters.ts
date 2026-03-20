import type { Character } from "../../packages/types/src/sim.js";

export const mockHutao: Character = {
  name: "hutao",
  level: 90,
  element: "pyro",
  max_level: 90,
  cons: 1,
  weapon: {
    name: "staffofhoma",
    refine: 1,
    level: 90,
    max_level: 90,
  },
  talents: {
    attack: 10,
    skill: 10,
    burst: 10,
  },
  stats: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
  snapshot: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
  sets: {
    crimsonwitchofflames: 4,
  },
};

export const mockXingqiu: Character = {
  name: "xingqiu",
  level: 90,
  element: "hydro",
  max_level: 90,
  cons: 6,
  weapon: {
    name: "sacrificialsword",
    refine: 5,
    level: 90,
    max_level: 90,
  },
  talents: {
    attack: 1,
    skill: 10,
    burst: 13,
  },
  stats: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
  snapshot: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
  sets: {
    emblemofseveredfate: 4,
  },
};

export const mockCharacters: Character[] = [mockHutao, mockXingqiu];
