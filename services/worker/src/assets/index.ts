export async function handleAssets(
  request: Request,
  event: FetchEvent
): Promise<Response> {
  const cacheUrl = new URL(request.url);
  const cacheKey = new Request(cacheUrl.toString(), request);
  console.log(`checking for cache key: ${cacheUrl}`);
  const cache = caches.default;
  let response = await cache.match(cacheKey);

  if (!response) {
    console.log(
      `Response for request url: ${request.url} not present in cache. Fetching and caching request.`
    );

    ///api/assets/avatar/cyno.png
    const key = new URL(request.url).pathname.replace("/api/assets/", "");
    console.log(`getting ${key}`);

    const object = await GCSIM_ASSETS.get(key);

    if (object === null) {
      console.log(`${key} not found in r2`);
      return new Response(`Not Found`, { status: 404 });
    }

    const headers = new Headers();
    object.writeHttpMetadata(headers);
    headers.set("etag", object.httpEtag);
    headers.set("Cache-Control", "max-age=5184000");

    response = new Response(object.body, {
      headers,
    });

    event.waitUntil(cache.put(cacheKey, response.clone()));
  }

  return response;
}
