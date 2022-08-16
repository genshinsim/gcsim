import { dbClient } from "..";
import { userData } from "../share/user";
import JSONbig from 'json-bigint'

export async function createOrGetUser(
  id: BigInt,
  name: string
): Promise<userData> {
  //make request to postgrest manually so that we're properly encoding the id
  const url = POSTGREST_ENDPOINT + "/rpc/get_or_insert_user"
  const req = new Request(
    url,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSONbig.stringify({
        id,
        name
      }),
    }
  )
  const resp = await fetch(req)
  const data = await resp.json<any[]>()
  if (!resp.ok) {
    console.log(`status: ${resp.status}`)
    console.log(data)
    throw new Error(`Failed to fetch user data`)
  }

  //body should be json array

  if (data === null) {
    throw "Unexpected no data returned";
  }

  if (data.length < 1) {
    throw "Unexpected no result rows";
  }

  const rows: userData[] = data;
  return rows[0];
}
