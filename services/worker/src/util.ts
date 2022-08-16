import JSONbig from 'json-bigint';

export async function postgRESTFetch(
  endpoint: string,
  requestInitr?: Request | RequestInit | undefined
): Promise<any> {
  const resp = await fetch(POSTGREST_ENDPOINT + endpoint, requestInitr);

  const str = await resp.text();
  const data = JSONbig.parse(str);

  if (!resp.ok) {
    throw new Error(`request failed; status: ${resp.status}, body: ${str}`);
  }

  return data;
}
