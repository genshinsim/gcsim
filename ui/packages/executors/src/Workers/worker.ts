/* eslint-disable @typescript-eslint/no-explicit-any */
// @ts-ignore
self.importScripts("/wasm_exec.js");

// @ts-ignore
let go: Go;

// @ts-ignore
function ready(req: { module: WebAssembly.Module }) {
	go = new Go();
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
				fatal: true,
			});
		});
}

// Handles a request that calls into Go. The Go functions return their errors, so if the call
// throws, or the Go program exits during it (a fatal error such as running out of memory), this
// instance is unusable: the fatal response makes the executor replace the worker.
// @ts-ignore
function callGo(handle: () => any): any {
	let reason =
		"its Go program exited (the console shows why, e.g. out of memory)";
	try {
		const resp = handle();
		if (!go.exited) {
			return resp;
		}
	} catch (e) {
		console.error(e);
		if (!go.exited) {
			reason = `${e}`;
		}
	}
	return {
		type: WorkerResponse.Failed,
		reason: `A sim worker crashed: ${reason}`,
		fatal: true,
	};
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
	if (typeof resp === "string" || resp instanceof String) {
		return {
			type: WorkerResponse.Failed,
			reason: JSON.parse(resp as string).error,
		};
	}
	return { type: WorkerResponse.Done, result: resp, itr: req.itr };
}

// @ts-ignore
function handleRequest(req: any) {
	switch (req.type as WorkerRequest) {
		case WorkerRequest.Ready:
			return ready(req);
		case WorkerRequest.Initialize:
			return respond(
				req,
				callGo(() => initialize(req)),
			);
		case WorkerRequest.Run: {
			const resp = callGo(() => run(req));
			// transfer the result's buffer instead of copying it
			return respond(
				req,
				resp,
				resp.result instanceof Uint8Array ? [resp.result.buffer] : [],
			);
		}
		default:
			console.error("aggregator - unknown request: ", req);
			throw new Error("aggregator unknown request");
	}
}
self.onmessage = (ev) => handleRequest(ev.data);

// Echoes the request's run id so the executor can drop responses from a cancelled run.
// @ts-ignore
function respond(
	req: { run: number },
	resp: object,
	transfer: Transferable[] = [],
) {
	postMessage({ ...resp, run: req.run }, transfer);
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
