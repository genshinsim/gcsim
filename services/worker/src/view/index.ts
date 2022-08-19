import { Request } from 'itty-router';
import { dbClient } from '..';

export async function handleView(
  request: Request,
  event: FetchEvent
): Promise<Response> {
  let { params } = request;
  if (!params || !params.key) {
    return new Response(null, {
      status: 400,
      statusText: 'Bad Request',
    });
  }

  const key = params.key;

  if (key === '') {
    return new Response(null, {
      status: 400,
      statusText: 'Bad Request',
    });
  }

  console.log(key);

  const cacheUrl = new URL(request.url);
  const cacheKey = new Request(cacheUrl.toString(), request);
  console.log(`checking for cache key: ${cacheUrl}`);
  const cache = caches.default;

  let response = await cache.match(cacheKey);

  if (!response) {
    console.log(
      `Response for request url: ${request.url} not present in cache. Fetching and caching request.`
    );

    const { data, error } = await dbClient
      .from('active_sim')
      .select()
      .eq('simulation_key', key);

    if (error !== null) {
      console.log(error);
      return new Response(null, {
        status: 500,
        statusText: 'Internal Server Error',
      });
    }

    if (data.length === 0) {
      return new Response('invalid key', {
        status: 400,
        statusText: 'Bad Request',
      });
    }

    response = new Response(
      JSON.stringify({
        data: data[0].viewer_file,
        meta: data[0].metadata,
      }),
      {
        headers: {
          'content-type': 'application/json',
          'Content-Encoding': 'gzip',
          'Cache-Control': 's-maxage=5184000',
        },
      }
    );

    event.waitUntil(cache.put(cacheKey, response.clone()));
  }

  return response;
}
