import { Request } from 'itty-router';
import { dbClient } from '..';

export async function handleListDBSims(request: Request): Promise<Response> {
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

  const { data, error } = await dbClient
    .from('db_sims_by_avatar')
    .select()
    .eq('avatar_name', key);

  if (error !== null) {
    console.log(error);
    return new Response(null, {
      status: 500,
      statusText: 'Internal Server Error',
    });
  }
  return new Response(JSON.stringify(data), {
    headers: {
      'content-type': 'application/json',
      'Content-Encoding': 'gzip',
    },
  });
}
