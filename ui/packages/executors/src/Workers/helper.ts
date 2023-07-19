/* eslint-disable @typescript-eslint/no-explicit-any */
// @ts-ignore
self.importScripts("/wasm_exec.js");

if (!WebAssembly.instantiateStreaming) {
  // polyfill
  WebAssembly.instantiateStreaming = async (resp, importObject) => {
    const source = await (await resp).arrayBuffer();
    return await WebAssembly.instantiate(source, importObject);
  };
}

let readyState = false;

// @ts-ignore
function ready(req: { wasm: string }) {
  const go = new Go();
  WebAssembly.instantiateStreaming(fetch(req.wasm), go.importObject)
      .then((result) => {
        go.run(result.instance);
        console.log("helper loaded okay");
        readyState = true;
      }).catch((e) => {
        console.error(e);
        postMessage({
            type: HelpResponse.Failed,
            reason: e instanceof Error ? e.message : "Unknown Error" });
      });
}

function validate(req: { id: number, cfg: string }) {
  const resp = JSON.parse(validateConfig(req.cfg));
  if (resp.error) {
    return { type: HelpResponse.Failed, reason: resp.error, id: req.id };
  }
  return { type: HelpResponse.Validate, cfg: resp, id: req.id };
}

function doSample(req: { id: number, cfg: string, seed: string }) {
  const resp = JSON.parse(sample(req.cfg, req.seed));
  if (resp.error) {
    return { type: HelpResponse.Failed, reason: resp.error, id: req.id };
  }
  return { type: HelpResponse.Sample, sample: resp, id: req.id };
}

// @ts-ignore
function handleRequest(req: any): any {
  switch (req.type as HelpRequest) {
    case HelpRequest.Validate:
      return validate(req);
    case HelpRequest.Sample:
      return doSample(req);
    default:
      console.error("helper - unknown request: ", req);
      throw new Error("helper unknown request");
  }
}

const queue: MessageEvent<any>[] = [];
onmessage = (ev) => {
  if (ev.data.type == HelpRequest.Ready) {
    ready(ev.data);
    return;
  }

  queue.push(ev);
  tryProcess();
};

function tryProcess() {
  if (!readyState) {
    setTimeout(tryProcess, 100);
    return;
  }

  const event = queue.shift();
  if (event) {
    postMessage(handleRequest(event.data));
  }
}

// TODO: I hate this
// Web Workers do not currently support modules (in all browsers), so instead all the relevant code in common
// has to be copy/pasted over
// Clean up when supported: https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Modules

enum HelpRequest {
  Ready = "ready",
  Validate = "validate",
  Sample = "sample",
}

enum HelpResponse {
  Failed = "failed",
  Validate = "validated",
  Sample = "sample",
}