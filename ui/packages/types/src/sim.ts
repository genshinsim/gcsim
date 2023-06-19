export interface SimResults {
  schema_version?: Version;
  sim_version?: string;
  version?: string; // legacy only
  build_date?: string;
  modified?: boolean;
  
  initial_character?: string;
  character_details?: Character[];
  target_details?: Enemy[];
  simulator_settings?: Settings;
  player_position?: Coord;

  config_file?: string;
  sample_seed?: string;

  statistics?: Statistics;
}

export interface Sample {
  config?: string;
  initial_character?: string;
  character_details?: Character[];
  target_details?: unknown[];
  seed?: string;
  logs?: LogDetails[];
}

export interface Version {
  major: string;
  minor: string;
}

export interface Statistics {
  // metadata
  min_seed?: string;
  max_seed?: string;
  p25_seed?: string;
  p50_seed?: string;
  p75_seed?: string;
  runtime?: number;
  iterations?: number;

  // summary
  duration?: SummaryStat;
  dps?: SummaryStat;
  rps?: SummaryStat;
  eps?: SummaryStat;
  hps?: SummaryStat;
  shp?: SummaryStat;

  // warnings
  warnings?: Warnings;

  // character stats
  failed_actions?: FailedActions[];

  character_dps?: FloatStat[];
  target_dps?: TargetDPS;
  element_dps?: ElementDPS;
  dps_by_element?: ElementStats[];
  dps_by_target?: TargetStats[];

  damage_buckets?: BucketStats;
  cumu_damage_contrib?: CharacterBucketStats;

  shields?: Shields;
}

export interface SummaryStat {
  min?: number;
  max?: number;
  mean?: number;
  sd?: number;
  q1?: number;
  q2?: number;
  q3?: number;
  histogram?: number[];
}

export interface Warnings {
  target_overlap?: boolean;
  insufficient_energy?: boolean;
  insufficient_stamina?: boolean;
  swap_cd?: boolean;
  skill_cd?: boolean;
}

export interface FailedActions {
  insufficient_energy?: FloatStat;
  insufficient_stamina?: FloatStat;
  swap_cd?: FloatStat;
  skill_cd?: FloatStat;
}

export interface Shields {
  [key: string]: ShieldInfo;
}

export interface ShieldInfo {
  hp?: Map<string, FloatStat>;
  uptime?: FloatStat;
}

export interface ElementStats {
  elements?: ElementDPS
}

export interface ElementDPS {
  [key: string]: FloatStat;
}

export interface TargetStats {
  targets?: TargetDPS
}

export interface TargetDPS {
  [key: string]: FloatStat;
}

export interface CharacterBucketStats {
  bucket_size?: number;
  characters?: CharacterBuckets[];
}

export interface CharacterBuckets {
  buckets: FloatStat[];
}

export interface BucketStats {
  bucket_size?: number;
  buckets?: FloatStat[];
}

export interface FloatStat {
  min?: number;
  max?: number;
  mean?: number;
  sd?: number;
}

export interface Settings {
  iterations?: number;
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

export interface Enemy {
  level?: number;
  hp?: number;
  resist?: { [key: string]: number };
  position?: Coord;
  particle_drop_threshold?: number;
  particle_drop_count?: number;
  particle_element?: number;
}

export interface Coord {
  x: number;
  y: number;
  r: number;
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
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  logs: { [key in string]: any };
  ordering?: { [key: string]: number };
};

export type StatusType = "idle" | "loading" | "done" | "error";

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
