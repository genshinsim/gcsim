import { IRequest } from 'itty-router';

export async function handleEnka(request: IRequest): Promise<Response> {
  let { params } = request;
  if (!params || !params.key) {
    return new Response(null, {
      status: 400,
      statusText: 'Bad Request',
    });
  }

  const uid = params.key;
  if (!/([1,2,5-9])\d{8}/.test(uid)) {
    return new Response(null, {
      status: 400,
      statusText: 'Bad Request',
    });
  }
  console.log(uid);

  const init = {
    headers: {
      "User-Agent": "gcsim/1.0",
    },
  };

  const resp = await fetch(`https://enka.network/api/uid/${uid}`, init);
  const contentType = resp.headers.get('content-type') || '';
  if (!resp.ok || !contentType.includes('application/json')) {
    return new Response(null, {
      status: 400,
      statusText: 'Bad Request',
    });
  }

  let avatars = [];
  const d = await resp.json();
  for (let avatar of d.avatarInfoList) {
    avatars.push(avatar);
  }

  if (d.owner != undefined) {
    console.log(uid, "has enka profile");

    const resp = await fetch(`https://enka.network/api/profile/${d.owner.username}/hoyos/${d.owner.hash}/builds/`, init);
    const contentType = resp.headers.get('content-type') || '';
    if (resp.ok && contentType.includes('application/json')) {
      const d = await resp.json();
      for (let [ _, builds ] of Object.entries(d)) {
        for (let build of builds) {
          if (build.live) { continue }
          build.avatar_data.name = build.name;
          avatars.push(build.avatar_data);
        }
      }
    }
  }

  return new Response(JSON.stringify(avatars), {
    status: resp.status,
    statusText: resp.statusText,
    headers: {
      'Content-Type': 'application/json',
      'Content-Encoding': 'gzip',
    },
  });
}
