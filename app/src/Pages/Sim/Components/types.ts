import { SlotKey, ISubstat, StatKey } from "./goodTypes";
export interface Weapon {
  key: string;
  name: string;
  icon: string;
  level: number;
  ascension: number;
  refinement: number;
}

export interface Artifact {
  setKey: string;
  slotKey: SlotKey;
  icon: string;
  rarity: number;
  level: number;
  mainStatKey: StatKey | "";
  substats: ISubstat[];
}

export interface Character {
  key: string;
  name: string;
  element: string;
  icon: string;
  level: number;
  constellation: number;
  ascension: number;
  talent: { auto: number; skill: number; burst: number };
  weapontype: string;
  weapon: Weapon;
  artifact: {
    flower: Artifact;
    plume: Artifact;
    sands: Artifact;
    goblet: Artifact;
    circlet: Artifact;
  };
}
