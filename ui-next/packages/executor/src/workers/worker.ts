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
      postMessage({ type: WorkerResponse.Ready });
    })
    .catch((e) => {
      postMessage({
        type: WorkerResponse.Failed,
        reason: e instanceof Error ? e.message : "Unknown Error",
      });
    });
}

function initialize(req: { cfg: string }) {
  const resp = initializeWorker(req.cfg);
  if (resp != null) {
    return { type: WorkerResponse.Failed, reason: JSON.parse(resp).error };
  }
  return { type: WorkerResponse.Initialized };
}

function run(req: { itr: number }) {
  try {
    const resp = simulate();
    if (typeof resp === "string" || resp instanceof String) {
      return {
        type: WorkerResponse.Failed,
        reason: JSON.parse(resp as string).error,
      };
    }
    return { type: WorkerResponse.Done, result: resp, itr: req.itr };
  } catch (e) {
    return { type: WorkerResponse.Failed, reason: `Failed with error: ${e}` };
  }
}

// biome-ignore lint/suspicious/noExplicitAny: worker message handler needs flexible typing
function handleRequest(req: any) {
  switch (req.type as WorkerRequest) {
    case WorkerRequest.Ready:
      return ready(req);
    case WorkerRequest.Initialize:
      return postMessage(initialize(req));
    case WorkerRequest.Run:
      return postMessage(run(req));
    default:
      throw new Error("worker unknown request");
  }
}
self.onmessage = (ev) => handleRequest(ev.data);

// Duplicated enums from common.ts because Web Workers cannot use ES module imports
// in all browsers. These must stay in sync with the common.ts definitions.

enum WorkerRequest {
  Ready = "ready",
  Initialize = "initialize",
  Run = "run",
}

enum WorkerResponse {
  Failed = "failed",
  Ready = "ready",
  Initialized = "initialized",
  Done = "done",
}
