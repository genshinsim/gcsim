/// <reference path="./wasm-types.d.ts" />

// Worker file — runs in a dedicated Web Worker context.
// Cannot use ES module imports since workers may not support modules in all browsers.
// Message types are duplicated from common.ts as enums (see note at bottom).

self.importScripts("/wasm_exec.js");

if (!WebAssembly.instantiateStreaming) {
  WebAssembly.instantiateStreaming = async (resp, importObject) => {
    const source = await (await resp).arrayBuffer();
    return await WebAssembly.instantiate(source, importObject);
  };
}

let readyState = false;

function ready(req: { wasm: string }) {
  const go = new Go();
  WebAssembly.instantiateStreaming(fetch(req.wasm), go.importObject)
    .then((result) => {
      go.run(result.instance);
      readyState = true;
    })
    .catch((e) => {
      postMessage({
        type: HelpResponse.Failed,
        reason: e instanceof Error ? e.message : "Unknown Error",
      });
    });
}

function validate(req: { id: number; cfg: string }) {
  const resp = JSON.parse(validateConfig(req.cfg));
  if (resp.error) {
    return { type: HelpResponse.Failed, reason: resp.error, id: req.id };
  }
  return { type: HelpResponse.Validate, cfg: resp, id: req.id };
}

function doSample(req: { id: number; cfg: string; seed: string }) {
  const resp = JSON.parse(sample(req.cfg, req.seed));
  if (resp.error) {
    return { type: HelpResponse.Failed, reason: resp.error, id: req.id };
  }
  return { type: HelpResponse.Sample, sample: resp, id: req.id };
}

// biome-ignore lint/suspicious/noExplicitAny: worker message handler needs flexible typing
function handleRequest(req: any): any {
  switch (req.type as HelpRequest) {
    case HelpRequest.Validate:
      return validate(req);
    case HelpRequest.Sample:
      return doSample(req);
    default:
      throw new Error("helper unknown request");
  }
}

// biome-ignore lint/suspicious/noExplicitAny: worker message queue holds arbitrary events
const queue: MessageEvent<any>[] = [];
self.onmessage = (ev) => {
  if (ev.data.type === HelpRequest.Ready) {
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

// Duplicated enums from common.ts because Web Workers cannot use ES module imports
// in all browsers. These must stay in sync with the common.ts definitions.

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
