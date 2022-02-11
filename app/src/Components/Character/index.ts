import { WeaponDetail } from "/src/Components/Weapon";

export * from "./CharacterCard";
export * from "./Stats";
export * from "./CharacterEdit";
export * from "./CharacterEditStats";
export * from "./CharacterEditWeapon";

export interface CharDetail {
  name: string;
  level: number;
  element: string;
  max_level: number;
  cons: number;
  weapon: WeaponDetail;
  talents: {
    attack: number;
    skill: number;
    burst: number;
  };
  stats: number[];
  snapshot: number[];
  sets: { [key: string]: number };
}
