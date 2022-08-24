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

//TODO: need to add create_time to this
export interface UserSimDetails {
  simulation_key: string;
  metadata: Metadata;
  is_permanent: boolean;
  user_id: BigInt;
  user_name: string;
}
