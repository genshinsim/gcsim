export type discordToken = {
  access_token: string;
  token_type: string;
  expires_in: number;
  refresh_token: string;
  scope: string;
};

export type discordUser = {
  id: string;
  username: string;
  discriminator: string;
};

export async function discordAccessToken(code: string): Promise<Response> {
  //grab token from discord
  let options = {
    url: 'https://discord.com/api/oauth2/token',
    method: 'POST',
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded;charset=utf-8',
    },
    body: new URLSearchParams({
      client_id: DISCORD_ID,
      client_secret: DISCORD_SECRET,
      grant_type: 'authorization_code',
      code: code,
      redirect_uri: `https://gcsim.app/auth/discord`,
      scope: 'identify',
    }).toString(),
  };
  console.log(options.body.toString());
  return fetch('https://discord.com/api/oauth2/token', options);
}

export async function requestDiscordId(
  bearer: string,
  token: string
): Promise<Response> {
  //grab identity
  const reqIdentity = new Request(`https://discord.com/api/users/@me`, {
    method: 'GET',
    headers: {
      authorization: `${bearer} ${token}`,
    },
  });
  return fetch(reqIdentity);
}
