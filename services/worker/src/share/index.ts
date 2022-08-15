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

  //key is uuid but no -
  // let key = uuid();
  // console.log(key);
  let perm = false;

  //check if this is a logged in user; if not then it can't be perm
  let user: userData | null = null;

  let id = await verifyToken(request.headers.get('X-AUTH-TOKEN'));
  if (id !== null) {
    user = await getUserInfo(id);
  }

  if (content.perm && user !== null) {
    perm = user.count < userLimits(user.user_role);
  }

  //store it
  const { data, error } = await dbClient.rpc('share_sim', {
    metadata: JSON.stringify(content.meta),
    viewer_file: content.data,
    user_id: user ? user.user_id : null,
    is_permanent: perm,
    is_public: false,
  });
  if (error !== null) {
    console.log(error);
    return new Response(null, {
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
      console.log(error);
      return new Response(null, {
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
  if (token !== null) {
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
