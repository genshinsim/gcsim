import { IRequest } from "itty-router";

export async function handleWasm(
  request: IRequest,
  event: FetchEvent
): Promise<Response> {
  const cacheUrl = new URL(request.url);
  const cacheKey = new Request(cacheUrl.toString(), request);
  console.log(`checking for cache key: ${cacheUrl}`);
  const cache = caches.default;
  let response = await cache.match(cacheKey);

  if (!response) {
    ///api/wasm/<branch>/<hash>/main.wasm
    const key = new URL(request.url).pathname.replace("/api/wasm/", "");
    console.log(`request key ${key}`);

    const object = await GCSIM_WASM.get(key);

    if (object == null) {
      console.log(`${key} not found in r2`);
      return new Response(`Not Found`, { status: 404 });
    }

    const headers = new Headers();
    object.writeHttpMetadata(headers);
    headers.set("etag", object.httpEtag);
    headers.set("Cache-Control", "max-age=5184000");
    headers.set("Content-Type", "application/wasm");

    response = new Response(object.body, {
      headers,
    });

    event.waitUntil(cache.put(cacheKey, response.clone()));
  }

  return response
}
