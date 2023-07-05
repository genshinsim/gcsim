import { IRequest } from "itty-router";
import { validator } from "./validation";

export async function handleShare(request: IRequest): Promise<Response> {
  let content: any;
  console.log("share request received! processing data");
  try {
    content = await request.json();
  } catch {
    return new Response(null, {
      status: 400,
      statusText: "Bad Request (Invalid JSON)",
    });
  }

  //validate input
  const valid = validator.validate(content);

  if (!valid.valid) {
    console.log(valid.errors);
    return new Response(null, { status: 400, statusText: "Bad Request" });
  }

  //post to endpoint

  return fetch(new Request(API_ENDPOINT + "/api/share"), {
    method: "POST",
    body: JSON.stringify(content),
    headers: request.headers,
  });
}
