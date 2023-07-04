// import { Request } from "itty-router";

export async function handleView(
  request: Request,
  event: FetchEvent
): Promise<Response> {
  let { params } = request;
  if (!params || !params.key) {
    return new Response(null, {
      status: 400,
      statusText: "Bad Request",
    });
  }

  const key = params.key;

  if (key === "") {
    return new Response(null, {
      status: 400,
      statusText: "Bad Request",
    });
  }

  console.log(key);

  const cacheUrl = new URL(request.url);
  const cacheKey = new Request(cacheUrl.toString(), request);
  console.log(`checking for cache key: ${cacheUrl}`);
  const cache = caches.default;

  //check if this is db route
  let dbStr = ""
  if (request.url.includes("/db/")) {
    dbStr = "db/"
  }

  let response = await cache.match(cacheKey);

  if (!response) {
    console.log(
      `Response for request url: ${request.url} not present in cache. Fetching and caching request.`
    );

    response = await fetch(new Request(API_ENDPOINT + "/api/share/" + dbStr + key));

    response = new Response(response.body, response);
    response.headers.append("Cache-Control", "s-maxage=1800");
    response.headers.append("Content-Encoding", "gzip");

    event.waitUntil(cache.put(cacheKey, response.clone()));
  } else {
    console.log(`cache hit for: ${request.url}`);
  }

  return response;
}
