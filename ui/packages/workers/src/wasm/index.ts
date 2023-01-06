export async function handleWasm(request: Request, event: FetchEvent): Promise<Response> {
  ///api/wasm/<branch>/<hash>/main.wasm
  const key = new URL(request.url).pathname.replace("/api/wasm/", "");
  console.log(`request key ${key}`);

  const object = await GCSIM_WASM.get(key);

  if (object == null) {
    console.log(`${key} not found in r2`);
    return new Response(`Not Found`, { status: 404 });
  }

  const headers = new Headers();
  object.writeHttpMetadata(headers);
  headers.set("etag", object.httpEtag);
  headers.set("Cache-Control", "max-age=5184000");
  headers.set("Content-Type", "application/wasm");

  return new Response(object.body, { headers });
}