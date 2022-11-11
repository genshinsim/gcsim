export async function handleAssets(
  request: Request,
  event: FetchEvent
): Promise<Response> {
  ///api/assets/avatar/cyno.png
  const key = new URL(request.url).pathname.replace("/api/assets/", "");
  console.log(`getting ${key}`);

  const object = await GCSIM_ASSETS.get(key);

  if (object === null) {
    console.log(`${key} not found in r2`);
    return new Response(`Not Found`, { status: 404 });
  }

  const headers = new Headers();
  object.writeHttpMetadata(headers);
  headers.set("etag", object.httpEtag);
  headers.set("Cache-Control", "max-age=5184000");

  return new Response(object.body, {
    headers,
  });
}
