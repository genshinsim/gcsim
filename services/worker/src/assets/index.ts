export async function handleAssets(
  request: Request,
  event: FetchEvent
): Promise<Response> {
  const cacheUrl = new URL(request.url).pathname;
  const endpoint = ASSETS_ENDPOINT + cacheUrl;
  const cacheKey = new Request(endpoint, request);
  console.log(`checking for cache key: ${endpoint}`);
  const cache = caches.default;

  let response = await cache.match(cacheKey);

  if (!response) {
    console.log(
      `Response for request url: ${endpoint} not present in cache. Fetching and caching request.`
    );

    const resp = await fetch(new Request(endpoint), {
      cf: {
        cacheTtl: 60 * 24 * 60 * 60,
        cacheEverything: true,
      },
    });

    response = new Response(resp.body, resp);
    response.headers.set('Cache-Control', 'max-age=5184000');

    event.waitUntil(cache.put(cacheKey, response.clone()));
  }

  return response;
}
