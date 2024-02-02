import { IRequest } from "itty-router";

export async function proxyRequest(request: IRequest): Promise<Response> {
  const x = new URL(request.url);
  return fetch(new Request(API_ENDPOINT + x.pathname + x.search, request));
}
