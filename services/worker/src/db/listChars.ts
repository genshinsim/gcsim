import { dbClient } from '..';

export async function handleListDBChars(request: Request): Promise<Response> {
  const { data, error } = await dbClient.from('avatars_with_db_sims').select();

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
