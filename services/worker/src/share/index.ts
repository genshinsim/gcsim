import { uuid } from '@cfworker/uuid';
import jwt from '@tsndr/cloudflare-worker-jwt';
import { dbClient } from '..';
import { getUserInfo, userData, userLimits } from './user';
import { uploadData, validator } from './validation';

/**
 * handles viewer share request
 * @param request incoming http request
 */
export async function handleShare(request: Request): Promise<Response> {
  // check content type
  let content: uploadData;
  console.log('share request received! processing data');
  try {
    content = await request.json<uploadData>();
  } catch {
    return new Response(null, {
      status: 400,
      statusText: 'Bad Request (Invalid JSON)',
    });
  }

  //validate input
  const valid = validator.validate(content);

  if (!valid.valid) {
    console.log(valid.errors);
    return new Response(null, { status: 400, statusText: 'Bad Request' });
  }

  // console.log('share received: ', content);

  //TODO: everything is perm for now
  let perm = true;

  //check if this is a logged in user; if not then it can't be perm
  // let user: userData | null = null;

  // let id = await verifyToken(request.headers.get('X-AUTH-TOKEN'));
  // console.log('user id: ', id);
  // if (id !== null) {
  //   try {
  //     user = await getUserInfo(id);
  //     console.log('got user info: ', user.user_id);
  //   } catch (error) {
  //     return new Response(JSON.stringify(error), {
  //       status: 500,
  //       statusText: 'Internal Server Error',
  //     });
  //   }
  // }

  // if (content.perm && user !== null) {
  //   perm = user.count < userLimits(user.user_role);
  //   console.log('user perm check: ', perm, user.user_id);
  // }

  //store it
  const { data, error } = await dbClient.rpc('share_sim', {
    metadata: JSON.stringify(content.meta),
    viewer_file: content.data,
    // user_id: user ? user.user_id : null,
    user_id: null,
    is_permanent: perm,
    is_public: false,
  });
  if (error !== null) {
    return new Response(JSON.stringify(error), {
      status: 500,
      statusText: 'Internal Server Error',
    });
  }

  //data is expected to be the key
  const key = data;

  //TODO: avatar and embed we don't actually care about the return value
  //and if they error'd or not

  //upload avatar information
  for (const char of content.meta.char_names) {
    console.log('linking character: ', char);
    const { error } = await dbClient.rpc('link_avatar_to_sim', {
      avatar: char,
      key: key,
    });
    if (error !== null) {
      return new Response(JSON.stringify(error), {
        status: 500,
        statusText: 'Internal Server Error',
      });
    }
  }

  //send request to generate embed
  await fetch(new Request(PREVIEW_ENDPOINT + '/embed/' + key), {
    method: 'POST',
  });

  if (error !== null) {
    console.log(error);
    return new Response(null, {
      status: 500,
      statusText: 'Internal Server Error',
    });
  }

  //TODO: cache in kv for 30 days

  return new Response(JSON.stringify({ key: key, perm: perm }), {
    headers: {
      'content-type': 'application/json',
    },
  });
}

/**
 * verifyToken takes a user supplied token via header and verifies if it is valid
 * @param token string representing the token; can be null if user did not supply a token
 * @returns user discord id if stored and is valid; otherwise null
 */
async function verifyToken(token: string | null): Promise<string | null> {
  console.log('verifying token ', token);
  if (token !== null && token !== '') {
    try {
      const ok = await jwt.verify(token, JWT_SECRET);
      if (ok) {
        const decoded = jwt.decode(token);
        if ('id' in decoded.payload) {
          return decoded.payload['id'];
        }
      }
    } catch (e) {
      //invalid token, do nothing
      console.log(e);
    }
  }
  return null;
}
