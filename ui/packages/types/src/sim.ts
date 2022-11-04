export interface SimResults {
  schema_version?: Version
  max_iterations?: number

  initial_character?: string
  character_details?: CharacterDetail[]

  config_file?: string
  debug?: LogDetails[]
  debug_seed?: string

  statistics?: Statistics
}

export interface Version {
  major: number
  minor: number
}

export interface CharacterDetail {
  name: string
}

export interface Statistics {
  // metadata
  min_seed?: string
  max_seed?: string
  p25_seed?: string
  p50_seed?: string
  p75_seed?: string
  runtime?: number
  iterations?: number

  // summary
  duration?: SummaryStat
  dps?: SummaryStat
  rps?: SummaryStat
  eps?: SummaryStat
  hps?: SummaryStat
  sps?: SummaryStat

  // warnings
  warnings?: Warnings

  // character stats
  failed_actions?: FailedActions[]
}

export interface SummaryStat {
  min?: number
  max?: number
  mean?: number
  sd?: number
  q1?: number
  q2?: number
  q3?: number
}

export interface Warnings {
  target_overlap?: boolean
  insufficient_energy?: boolean
  insufficient_stamina?: boolean
  swap_cd?: boolean
  skill_cd?: boolean
}

export interface FailedActions {
  insufficient_energy?: FloatStat
  insufficient_stamina?: FloatStat
  swap_cd?: FloatStat
  skill_cd?: FloatStat
}

export interface FloatStat {
  min?: number
  max?: number
  mean?: number
  sd?: number
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

export type LogDetails = {
  char_index: number;
  ended: number;
  event: string;
  frame: number;
  msg: string;
  logs: { [key in string]: unknown };
  ordering?: { [key: string]: number };
};

export type StatusType = 'idle' | 'loading' | 'done' | 'error';

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