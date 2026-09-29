/* eslint-disable @typescript-eslint/no-explicit-any */
// @ts-ignore
self.importScripts("/wasm_exec.js");

// @ts-ignore
function ready(req: { module: WebAssembly.Module }) {
	const go = new Go();
	WebAssembly.instantiate(req.module, go.importObject)
		.then((instance) => {
			go.run(instance);
			postMessage({ type: WorkerResponse.Ready });
		})
		.catch((e) => {
			console.error(e);
			postMessage({
				type: WorkerResponse.Failed,
				reason: e instanceof Error ? e.message : "Unknown Error",
			});
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
		console.log("simulate() call failed");
		return { type: WorkerResponse.Failed, reason: `Failed with error: ${e}` };
	}
}

// @ts-ignore
function handleRequest(req: any) {
	switch (req.type as WorkerRequest) {
		case WorkerRequest.Ready:
			return ready(req);
		case WorkerRequest.Initialize:
			return respond(req, initialize(req));
		case WorkerRequest.Run:
			return respond(req, run(req));
		default:
			console.error("aggregator - unknown request: ", req);
			throw new Error("aggregator unknown request");
	}
}
self.onmessage = (ev) => handleRequest(ev.data);

// Echoes the request's run id so the executor can drop responses from a cancelled run.
// @ts-ignore
function respond(req: { run: number }, resp: object) {
	postMessage({ ...resp, run: req.run });
}

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
