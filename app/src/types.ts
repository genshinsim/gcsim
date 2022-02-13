export interface Character {
  name: string;
  level: number;
  element: string;
  max_level: number;
  cons: number;
  weapon: Weapon;
  talents: Talent;
  stats: number[];
  snapshot: number[];
  sets: Set;
}

export const defaultStats = [
  0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
];

export const maxStatLength = defaultStats.length;

export interface Talent {
  attack: number;
  skill: number;
  burst: number;
}

export interface Set {
  [key: string]: number;
}

export interface Weapon {
  name: string;
  refine: number;
  level: number;
  max_level: number;
}

export interface Result {
  is_damage_mode: boolean;
  char_names: string[];
  char_details: Character[];
  damage_by_char: { [key in string]: number }[];
  damage_instances_by_char: { [key in string]: number }[];
  damage_by_char_by_targets: { [key in number]: number }[];
  char_active_time: number[];
  abil_usage_count_by_char: { [key in string]: number }[];
  particle_count: { [key in string]: number };
  reactions_triggered: { [key in string]: number };
  sim_duration: number;
  ele_uptime: { [key in number]: number }[];
  energy_when_burst: number[][];
  damage: number;
  dps: number;
  seed: number;
}
