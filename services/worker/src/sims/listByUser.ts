import { Request } from 'itty-router';
import { postgRESTFetch } from '../util';

export async function handleListUserSims(request: Request): Promise<Response> {
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

  const data = await postgRESTFetch(
    `/active_user_simulations?user_id=eq.${key}`
  );

  return new Response(JSON.stringify(data), {
    headers: {
      'content-type': 'application/json',
      'Content-Encoding': 'gzip',
    },
  });
}
