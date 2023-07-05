import { IRequest } from "itty-router";
import pako from "pako";

export async function handleLegacy(request: IRequest): Promise<Response> {
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

  let resp = await fetch(new Request(`https://gcsim.app/api/view/${key}`));

  if (resp.status != 200) {
    return resp;
  }

  try {
    let res: any = await resp.json();
    //we need to respond with res.
    const binaryStr = Uint8Array.from(atob(res.data), (v) => v.charCodeAt(0));
    const restored = pako.inflate(binaryStr, { to: "string" });

    let content = JSON.parse(restored);
    content.debug = "";
    content.text = "";
    // let keys = Object.keys(content);
    // console.log(keys);
    return new Response(JSON.stringify(content), {
      status: 200,
      statusText: "Ok",
    });
  } catch (err) {
    console.log(err);
    return new Response(null, {
      status: 500,
      statusText: "Error grabbing sim",
    });
  }
}
