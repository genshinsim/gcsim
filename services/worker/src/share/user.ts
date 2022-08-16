import { dbClient } from '..';
import { postgRESTFetch } from '../util';

/**
 * Returns the maximum number of perm share an user can have given the user's
 * role. Return -1 if unlimited
 * @param role user role value
 */
export function userLimits(role: number): number {
  if (role >= 99) {
    return -1;
  }
  //devs etc..
  if (role >= 20) {
    return 1000;
  }
  //subscribers
  if (role >= 10) {
    return 300;
  }
  //one time donation
  if (role >= 5) {
    return 100;
  }
  return 5;
}

export type userData = {
  user_id: BigInt;
  user_key: string;
  user_name: string;
  user_role: number;
  count: number;
};

export async function getUserInfo(id: string): Promise<userData> {
  try {
    const data = await postgRESTFetch(
      `/user_simulation_count?user_id=eq.${id}`
    );

    if (data === null) {
      throw 'Unexpected no data returned';
    }

    if (data.length < 1) {
      throw 'Unexpected no result rows';
    }

    return data[0];
  } catch (error) {
    throw error;
  }
}
