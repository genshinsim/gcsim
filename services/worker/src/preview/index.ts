const uuidRegex =
  /^[0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12}$/i;

export async function handlePreview(request: Request): Promise<Response> {
  //get last bit of url and change if valid uuid
  const segments = new URL(request.url).pathname.split("/");
  let last = segments.pop() || segments.pop(); // Handle potential trailing slash
  console.log(last);

  if (last === undefined) {
    return new Response(null, {
      status: 400,
      statusText: "Bad Request",
    });
  }
  //strip .png
  last = last.replace(".png", "");
  if (!last.match(uuidRegex)) {
    return new Response(null, {
      status: 400,
      statusText: "Bad Request",
    });
  }

  const resp = await fetch(new Request(PREVIEW_ENDPOINT + "/embed/" + last), {
    cf: {
      cacheTtl: 60 * 24 * 60 * 60,
      cacheEverything: true,
    },
  });

  const next = new Response(resp.body, resp);
  next.headers.set("Cache-Control", "max-age=5184000");
  return next;
}
