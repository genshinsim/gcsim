/* eslint-disable @typescript-eslint/no-explicit-any */
// @ts-ignore
self.importScripts("wasm_exec.js");

if (!WebAssembly.instantiateStreaming) {
  // polyfill
  WebAssembly.instantiateStreaming = async (resp, importObject) => {
    const source = await (await resp).arrayBuffer();
    return await WebAssembly.instantiate(source, importObject);
  };
}

// @ts-ignore
const go = new Go();
WebAssembly.instantiateStreaming(fetch('/main.wasm'), go.importObject)
  .then((result) => {
    go.run(result.instance);
    console.log("helper loaded okay");
  }).catch((e) => {
    console.error(e);
    postMessage({
      type: HelpResponse.Failed,
      reason: e instanceof Error ? e.message : "Unknown Error" });
  });

function validate(req: { cfg: string }) {
  const resp = JSON.parse(validateConfig(req.cfg));
  if (resp.error) {
    return { type: HelpResponse.Failed, reason: resp.error };
  }
  return { type: HelpResponse.Validate, cfg: resp };
}

function generateDebug(req: { cfg: string, seed: string }) {
  const resp = JSON.parse(debug(req.cfg, req.seed));
  if (resp.error) {
    return { type: HelpResponse.Failed, reason: resp.error };
  }
  return { type: HelpResponse.GenerateDebug, debug: resp };
}

// @ts-ignore
function handleRequest(req: any): any {
  switch (req.type as HelpRequest) {
    case HelpRequest.Validate:
      return validate(req);
    case HelpRequest.GenerateDebug:
      return generateDebug(req);
    default:
      console.error("helper - unknown request: ", req);
      throw new Error("helper unknown request");
  }
}
onmessage = (ev) => postMessage(handleRequest(ev.data));

// TODO: I hate this
// Web Workers do not currently support modules (in all browsers), so instead all the relevant code in common
// has to be copy/pasted over
// Clean up when supported: https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Modules

enum HelpRequest {
  Validate = "validate",
  GenerateDebug = "gen_debug",
}

enum HelpResponse {
  Failed = "failed",
  Validate = "validated",
  GenerateDebug = "gen_debug",
}