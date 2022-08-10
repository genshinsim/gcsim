import { dbClient } from "..";
import { userData } from "../share/user";

export async function createOrGetUser(
  key: string,
  name: string
): Promise<userData> {
  const { data, error } = await dbClient.rpc("get_or_insert_user", {
    key,
    name,
  });

  if (error !== null) {
    throw error;
  }

  if (data === null) {
    throw "Unexpected no data returned";
  }

  if (data.length < 1) {
    throw "Unexpected no result rows";
  }

  const rows: userData[] = data;
  return rows[0];
}
