const corsHeaders = {
  "Access-Control-Allow-Origin": "*",
  "Access-Control-Allow-Methods": "GET,HEAD,POST,OPTIONS",
  "Access-Control-Max-Age": "86400",
};

export async function handleOptions(request: Request): Promise<Response> {
  // Make sure the necessary headers are present
  // for this to be a valid pre-flight request
  let reqheader = request.headers;
  if (
    reqheader.get("Origin") !== null &&
    reqheader.get("Access-Control-Request-Method") !== null &&
    reqheader.get("Access-Control-Request-Headers") !== null
  ) {
    // Handle CORS pre-flight request.
    // If you want to check or reject the requested method + headers
    // you can do that here.
    let respHeaders: HeadersInit = {
      ...corsHeaders,
    };

    const acah = reqheader.get("Access-Control-Request-Headers");
    if (acah !== null) {
      // Allow all future content Request headers to go back to browser
      // such as Authorization (Bearer) or X-Client-Name-Version
      respHeaders["Access-Control-Allow-Headers"] = acah;
    }

    return new Response(null, {
      headers: respHeaders,
    });
  } else {
    // Handle standard OPTIONS request.
    // If you want to allow other HTTP Methods, you can do that here.
    return new Response(null, {
      headers: {
        Allow: "GET, HEAD, POST, OPTIONS",
      },
    });
  }
}
