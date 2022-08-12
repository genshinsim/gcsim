import { Request } from "itty-router";
import { dbClient } from "..";

export async function handleView(request: Request): Promise<Response> {
  let { params } = request;
  if (!params || !params.key) {
    return new Response(null, {
      status: 400,
      statusText: "Bad Request",
    });
  }

  const key = params.key;

  if (key === "") {
    return new Response(null, {
      status: 400,
      statusText: "Bad Request",
    });
  }

  console.log(key);

  //TODO: cache in kv here

  const { data, error } = await dbClient
    .from("active_simulations")
    .select()
    .eq("simulation_key", key);

  if (error !== null) {
    console.log(error);
    return new Response(null, {
      status: 500,
      statusText: "Internal Server Error",
    });
  }

  if (data.length === 0) {
    return new Response("invalid key", {
      status: 400,
      statusText: "Bad Request",
    });
  }

  return new Response(
    JSON.stringify({
      data: data[0].viewer_file,
      meta: data[0].metadata,
    }),
    {
      headers: {
        "content-type": "application/json",
        "Content-Encoding": "gzip",
      },
    }
  );
}
