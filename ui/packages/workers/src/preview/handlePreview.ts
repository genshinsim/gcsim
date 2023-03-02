export async function handlePreview(
  request: Request,
  event: FetchEvent
): Promise<Response> {
  //get last bit of url and change if valid uuid
  const pn = new URL(request.url).pathname
  const segments = pn.split("/");
  let last = segments.pop() || segments.pop(); // Handle potential trailing slash
  console.log(last);
  //if this ends in db/{key}, we need to stick in the db in front
  if (pn.includes("/db/")) {
    last = "db/" + last
  }

  if (last === undefined) {
    return new Response(null, {
      status: 400,
      statusText: "Bad Request",
    });
  }
  //strip .png
  last = last.replace(".png", "");

  const cacheUrl = new URL(request.url);
  const cacheKey = new Request(cacheUrl.toString(), request);
  console.log(`checking for cache key: ${cacheUrl}`);
  const cache = caches.default;
  let response = await cache.match(cacheKey);

  if (!response) {
    console.log(
      `Response for request url: ${request.url} not present in cache. Fetching and caching request.`
    );

    const resp = await fetch(
      new Request(API_ENDPOINT + "/preview/" + last),
      {
        cf: {
          cacheTtl: 60 * 24 * 60 * 60,
          cacheEverything: true,
        },
      }
    );

    response = new Response(resp.body, resp);
    response.headers.set("Cache-Control", "max-age=5184000");

    event.waitUntil(cache.put(cacheKey, response.clone()));
  }

  return response;
}
