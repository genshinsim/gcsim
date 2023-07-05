import { IRequest } from 'itty-router';

export async function handleEnka(request: IRequest): Promise<Response> {
  let { params } = request;
  if (!params || !params.key) {
    return new Response(null, {
      status: 400,
      statusText: 'Bad Request',
    });
  }

  const key = params.key;

  if (!/([1,2,5-9])\d{8}/.test(key)) {
    return new Response(null, {
      status: 400,
      statusText: 'Bad Request',
    });
  }

  console.log(key);

  const init = {
    headers: {
      Authorization: 'Bearer KEKWOMEGALULGCSIMEPICSRL',
    },
  };

  const requrl = `https://enka.network/api/uid/${key}/`;
  const response = await fetch(requrl, init);
  console.log(response);
  const results = await gatherResponse(response);
  const respHeaders = {
    'Content-Type': 'application/json',
    'Content-Encoding': 'gzip',
  };
  return new Response(results, {
    status: response.status,
    statusText: response.statusText,
    headers: respHeaders,
  });
}

/**
 * gatherResponse awaits and returns a response body as a string.
 * Use await gatherResponse(..) in an async function to get the response body
 * @param {Response} response
 */
async function gatherResponse(response: Response) {
  const { headers } = response;
  const contentType = headers.get('content-type') || '';
  if (contentType.includes('application/json')) {
    return JSON.stringify(await response.json());
  } else if (contentType.includes('application/text')) {
    return response.text();
  } else if (contentType.includes('text/html')) {
    return response.text();
  } else {
    return response.text();
  }
}
