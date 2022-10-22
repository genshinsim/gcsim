import { Character } from './sim';

export interface Metadata {
  char_names: string[];
  dps: SummaryStats;
  sim_duration: SummaryStats;
  dps_by_target: { [key: number]: SummaryStats };
  iter: number;
  runtime: number;
  num_targets: number;
  char_details: Character[];
}

export interface SummaryStats {
  mean: number;
  min: number;
  max: number;
  sd?: number;
}

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
  debug: string | [any];
  runtime: number;
  config_file: string;
  num_targets: number;
  //character details
  char_details: Character[];
}
