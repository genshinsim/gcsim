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

// @ts-ignore
function ready(req: { wasm: string }) {
  const go = new Go();
  WebAssembly.instantiateStreaming(fetch(req.wasm), go.importObject)
      .then((result) => {
        go.run(result.instance);
        postMessage({ type: WorkerResponse.Ready });
      }).catch((e) => {
        console.error(e);
        postMessage({
          type: WorkerResponse.Failed,
          reason: e instanceof Error ? e.message : "Unknown Error" });
      });
}

// @ts-ignore
function initialize(req: { cfg: string }) {
  const resp = initializeWorker(req.cfg);
  if (resp != null) {
    return { type: WorkerResponse.Failed, reason: JSON.parse(resp).error };
  }
  return { type: WorkerResponse.Initialized };
}

function run(req: { itr: number }) {
  const resp = simulate();
  if (typeof(resp) == "string" || resp instanceof String) {
    return { type: WorkerResponse.Failed, reason: JSON.parse(resp as string).error };
  }
  return { type: WorkerResponse.Done, result: resp, itr: req.itr };
}

// @ts-ignore
function handleRequest(req: any) {
  switch (req.type as WorkerRequest) {
    case WorkerRequest.Ready:
      return ready(req);
    case WorkerRequest.Initialize:
      return postMessage(initialize(req));
    case WorkerRequest.Run:
      return postMessage(run(req));
    default:
      console.error("aggregator - unknown request: ", req);
      throw new Error("aggregator unknown request");
  }
}
onmessage = (ev) => handleRequest(ev.data);

// TODO: I hate this
// Web Workers do not currently support modules (in all browsers), so instead the relevant code in common
// has to be copy/pasted over
// Clean up when supported: https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Modules

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