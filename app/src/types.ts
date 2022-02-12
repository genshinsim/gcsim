export interface Character {
  name: string;
  level: number;
  element: string;
  max_level: number;
  cons: number;
  weapon: Weapon;
  talents: {
    attack: number;
    skill: number;
    burst: number;
  };
  stats: number[];
  snapshot: number[];
  sets: { [key: string]: number };
}

export interface Weapon {
  name: string;
  refine: number;
  level: number;
  max_level: number;
}
