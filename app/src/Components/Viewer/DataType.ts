export interface SimResults {
  v2?: boolean;
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
  char_details: CharDetail[];
}

export interface CharDetail {
  name: string;
  level: number;
  element: string;
  max_level: number;
  cons: number;
  weapon: {
    name: string;
    refine: number;
    level: number;
    max_level: number;
  };
  talents: {
    attack: number;
    skill: number;
    burst: number;
  };
  stats: number[];
  snapshot: number[];
  sets: { [key: string]: number };
}

export interface SummaryStats {
  mean: number;
  min: number;
  max: number;
  sd?: number;
}
