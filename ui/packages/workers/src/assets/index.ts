import {IRequest} from 'itty-router';

export async function handleAssets(
  request: IRequest,
  event: FetchEvent,
): Promise<Response> {
  const cacheUrl = new URL(request.url);
  const cacheKey = new Request(cacheUrl.toString(), request);
  console.log(`checking for cache key: ${cacheUrl}`);
  const cache = caches.default;
  let response = await cache.match(cacheKey);

  if (!response) {
    console.log(
      `Response for request url: ${request.url} not present in cache. Fetching and caching request.`,
    );

    const resp = await fetch(new Request(ASSETS_ENDPOINT + cacheUrl.pathname), {
      cf: {
        cacheTtl: 60 * 24 * 60 * 60,
        cacheEverything: true,
      },
    });

    response = new Response(resp.body, resp);

    // only cache if response = 200 and this is not a placeholder
    if (
      resp.status === 200 &&
      resp.headers.get('Cache-Control') !== 'no-cache'
    ) {
      response.headers.set('Cache-Control', 'max-age=5184000');
      event.waitUntil(cache.put(cacheKey, response.clone()));
    }
  }

  return response;
}
