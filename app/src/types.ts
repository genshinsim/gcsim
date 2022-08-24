import { Character, Talent, Weapon, Set } from './Types/sim';

export const defaultStats = [
  0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
];

export const maxStatLength = defaultStats.length;

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

export interface SummaryStats {
  mean: number;
  min: number;
  max: number;
  sd?: number;
}

export interface DBItem {
  author: string;
  config: string;
  description: string;
  hash: string;
  team: DBCharInfo[];
  dps: number;
  mode: string;
  duration: number;
  target_count: number;
  viewer_key: string;
}

export interface DBCharInfo {
  name: string;
  con: number;
  weapon: string;
  refine: number;
  er: number;
  talents: Talent;
}

export interface ParsedResult {
  characters: ParsedCharacterProfile[];
  errors: string[];
  player_initial_pos: { x: number; y: number; r: number };
}
export interface ParsedCharacterProfile {
  base: Base;
  weapon: Weapon;
  talents: Talent;
  stats: number[];
  sets: Set;
}
export interface Base {
  key: string;
  name: string;
  element: string;
  level: number;
  max_level: number;
  base_hp: number;
  base_atk: number;
  base_def: number;
  cons: number;
  start_hp: number;
}
export interface ParsedWeapon {
  name: string;
  key: string;
  refine: number;
  level: number;
  max_level: number;
  base_atk: number;
}
