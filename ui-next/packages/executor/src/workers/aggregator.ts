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

function ready(req: { wasm: string }) {
  const go = new Go();
  WebAssembly.instantiateStreaming(fetch(req.wasm), go.importObject)
    .then((result) => {
      go.run(result.instance);
      postMessage({ type: AggResponse.Ready });
    })
    .catch((e) => {
      postMessage({
        type: AggResponse.Failed,
        reason: e instanceof Error ? e.message : "Unknown Error",
      });
    });
}

function initialize(req: { cfg: string }) {
  const resp = JSON.parse(initializeAggregator(req.cfg));
  if (resp.error) {
    return { type: AggResponse.Failed, reason: resp.error };
  }
  return { type: AggResponse.Initialized, result: resp };
}

function add(req: { result: Uint8Array }) {
  const resp = aggregate(req.result);
  if (resp != null) {
    return { type: AggResponse.Failed, reason: JSON.parse(resp).error };
  }
  return { type: AggResponse.Done };
}

function doFlush() {
  const resp = JSON.parse(flush());
  if (resp.error) {
    return { type: AggResponse.Failed, reason: resp.error };
  }
  return { type: AggResponse.Result, result: resp };
}

// biome-ignore lint/suspicious/noExplicitAny: worker message handler needs flexible typing
function handleRequest(req: any): any {
  switch (req.type as AggRequest) {
    case AggRequest.Ready:
      return ready(req);
    case AggRequest.Initialize:
      return postMessage(initialize(req));
    case AggRequest.Add:
      return postMessage(add(req));
    case AggRequest.Flush:
      return postMessage(doFlush());
    default:
      throw new Error("aggregator unknown request");
  }
}
self.onmessage = (ev) => handleRequest(ev.data);

// Duplicated enums from common.ts because Web Workers cannot use ES module imports
// in all browsers. These must stay in sync with the common.ts definitions.

enum AggRequest {
  Ready = "ready",
  Initialize = "initialize",
  Add = "add",
  Flush = "flush",
}

enum AggResponse {
  Failed = "failed",
  Ready = "ready",
  Initialized = "initialized",
  Done = "done",
  Result = "result",
}
