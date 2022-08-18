import {
  discordToken,
  discordAccessToken,
  requestDiscordId,
  discordUser,
} from './discord';
import jwt from '@tsndr/cloudflare-worker-jwt';
import { createOrGetUser } from './user';

export type jwtToken = {
  //standard stuff
  iss?: string;
  sub: string;
  aud?: string;
  exp: number;
  //app stuff
  name: string; //discord name
  id: string;
  role: number; //access lvl
};

const defHeaders: HeadersInit = {
  'content-type': 'application/json',
  'Content-Encoding': 'gzip',
};
const respFactory = (
  body: BodyInit | null,
  status: number,
  text: string
): Response =>
  new Response(body, { status: status, statusText: text, headers: defHeaders });

export async function handleAuth(request: Request): Promise<Response> {
  let req_headers = request.headers;

  //check for discord code
  const code = req_headers.get('X-DISCORD-CODE');
  if (code === null) {
    return respFactory(null, 403, 'Forbidden');
  }

  console.log('requesting for token: ', code);
  let token: discordToken;
  try {
    const resp = await discordAccessToken(code);

    token = await resp.json<discordToken>();
    console.log('response received: ', JSON.stringify(resp));
    console.log('response json: ', token);
    if (!resp.ok) {
      console.log('response not ok');
      return respFactory(null, resp.status, resp.statusText);
    }
    //check status
    if (resp.status !== 200) {
      console.log('response status not 200');
      return respFactory(null, resp.status, resp.statusText);
    }
  } catch (err) {
    console.log('response errored ', err);
    return respFactory(
      JSON.stringify({ err: err }),
      500,
      'Internal Server Error'
    );
  }
  console.log('token received: ', token);

  const resp = await requestDiscordId(token.token_type, token.access_token);

  if (!resp.ok) {
    console.log('error occured fetching from discord: ', resp.status);
    return respFactory(
      JSON.stringify({ err: resp.statusText }),
      500,
      'Internal Server Error'
    );
  }

  const discordIdentity: discordUser = await resp.json<discordUser>();
  console.log(discordIdentity);

  //generate a jwt containing the user's discord id
  const expiry = Math.floor(Date.now() / 1000) + 30 * 24 * 60 * 60; //30 days
  const userToken: jwtToken = {
    iss: 'gcsim.app',
    sub: 'user-token',
    exp: expiry,
    name: discordIdentity.username + '#' + discordIdentity.discriminator,
    id: discordIdentity.id,
    role: 0,
  };
  const secret = await jwt.sign(userToken, JWT_SECRET);

  //save the secret
  USER_TOKENS.put(discordIdentity.id, secret, { expiration: expiry });

  //get user info
  try {
    const id = BigInt(discordIdentity.id);
    const userData = await createOrGetUser(
      id,
      `${discordIdentity.username}#${discordIdentity.discriminator}`
    );
    return respFactory(
      JSON.stringify({ user: userData, token: secret }),
      200,
      ''
    );
  } catch (e) {
    console.log('error getting user: ', e);
    return respFactory(
      JSON.stringify({ err: e }),
      500,
      'Internal Server Error'
    );
  }
}
