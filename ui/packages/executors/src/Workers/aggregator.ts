/* eslint-disable @typescript-eslint/no-explicit-any */
// @ts-ignore
self.importScripts("/wasm_exec.js");

// @ts-ignore
function ready(req: { module: WebAssembly.Module }) {
	const go = new Go();
	WebAssembly.instantiate(req.module, go.importObject)
		.then((instance) => {
			go.run(instance);
			console.log("aggregator loaded okay");
			postMessage({ type: AggResponse.Ready });
		})
		.catch((e) => {
			console.error(e);
			postMessage({
				type: AggResponse.Failed,
				reason: e instanceof Error ? e.message : "Unknown Error",
			});
		});
}

// @ts-ignore
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
	// TODO: have a specific result response type to enforce (protos?)
	const resp = JSON.parse(flush());
	if (resp.error) {
		return { type: AggResponse.Failed, reason: resp.error };
	}
	return { type: AggResponse.Result, result: resp };
}

// @ts-ignore
function handleRequest(req: any): any {
	switch (req.type as AggRequest) {
		case AggRequest.Ready:
			return ready(req);
		case AggRequest.Initialize:
			return respond(req, initialize(req));
		case AggRequest.Add:
			return respond(req, add(req));
		case AggRequest.Flush:
			return respond(req, doFlush());
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
// Web Workers do not currently support modules (in all browsers), so instead all the relevant code in common
// has to be copy/pasted over
// Clean up when supported: https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Modules

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
