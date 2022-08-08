export interface ResultsSummary {
  v2?: boolean;
  version?: string;
  build_date?: string;
  is_damage_mode: boolean;
  active_char: string;
  char_names: string[];
  damage_by_char: { [key: string]: SummaryStats }[];
  damage_instances_by_char: { [key: string]: SummaryStats }[];
  damage_by_char_by_targets: { [key: number]: SummaryStats }[];
  char_active_time: SummaryStats[];
  abil_usage_count_by_char: { [key: string]: SummaryStats }[];
  particle_count: { [key: string]: SummaryStats };
  reactions_triggered: { [key: string]: SummaryStats };
  sim_duration: SummaryStats;
  ele_uptime: { [key: number]: SummaryStats }[];
  required_er: number[] | null;
  damage: SummaryStats;
  dps: SummaryStats;
  dps_by_target: { [key: number]: SummaryStats };
  damage_over_time: { [key: string]: SummaryStats };
  iter: number;
  text: string;
  debug: string;
  runtime: number;
  config_file: string;
  num_targets: number;
  //character details
  char_details: Character[];
}
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
  date_added?: string;
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
  player_initial_pos: { x: number, y: number, r: number};
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
