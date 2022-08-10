import { dbClient } from "..";

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
  user_id: number;
  user_key: string;
  user_name: string;
  user_role: number;
  count: number;
};

export async function getUserInfo(id: string): Promise<userData> {
  try {
    const res = await dbClient
      .from("user_simulation_count")
      .select()
      .eq("user_id", id);

    if (res.error === null) {
      const rows: userData[] = res.data;
      if (rows.length > 0) {
        return rows[0];
      }
    } else {
      console.log(res.error);
      throw res.error;
    }
  } catch (error) {
    throw error;
  }

  //unreachable
  throw "unreachable code";
}
