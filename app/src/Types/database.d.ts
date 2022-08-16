import { Metadata } from './stats';

export interface DBAvatarSimCount {
  avatar_name: string;
  sim_count: number;
}

export interface DBAvatarSimDetails {
  simulation_key: string;
  metadata: Metadata;
  is_permanent: boolean;
  create_time: number;
  git_hash: string;
  sim_description: string;
  config: string;
}
