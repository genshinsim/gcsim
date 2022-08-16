import { userData } from '../share/user';
import JSONbig from 'json-bigint';
import { postgRESTFetch } from '../util';

export async function createOrGetUser(
  id: BigInt,
  name: string
): Promise<userData> {
  //make request to postgrest manually so that we're properly encoding the id
  try {
    const data = await postgRESTFetch('/rpc/get_or_insert_user', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSONbig.stringify({
        id,
        name,
      }),
    });

    //body should be json array
    if (data === null) {
      throw 'Unexpected no data returned';
    }
    if (data.length < 1) {
      throw 'Unexpected no result rows';
    }

    const rows: userData[] = data;
    return rows[0];
  } catch (error) {
    throw error;
  }
}
