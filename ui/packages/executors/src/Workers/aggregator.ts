/* eslint-disable @typescript-eslint/no-explicit-any */
// @ts-ignore
self.importScripts("/wasm_exec.js");

// @ts-ignore
let go: Go;

// Asks the helper for the compiled module over req.port, worker to worker (see share in
// helper.ts).
// @ts-ignore
function ready(req: { port: MessagePort }) {
	req.port.onmessage = (ev) => {
		req.port.close();
		load(ev.data);
	};
	req.port.onmessageerror = () => {
		req.port.close();
		postMessage({
			type: AggResponse.Failed,
			reason: "The aggregator couldn't receive the compiled wasm module",
			fatal: true,
		});
	};
	req.port.postMessage(null);
}

// @ts-ignore
function load(module: WebAssembly.Module) {
	go = new Go();
	WebAssembly.instantiate(module, go.importObject)
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
				fatal: true,
			});
		});
}

// Handles a request that calls into Go. The Go functions return their errors, so if the call
// throws, or the Go program exits during it (a fatal error such as running out of memory), this
// instance is unusable: the fatal response makes the executor replace the aggregator.
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
		type: AggResponse.Failed,
		reason: `The aggregator crashed: ${reason}`,
		fatal: true,
	};
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

function doFlush(req: { final: boolean }) {
	const start = performance.now();
	const resp = flush();
	if (typeof resp === "string") {
		return { type: AggResponse.Failed, reason: JSON.parse(resp).error };
	}
	return {
		type: AggResponse.Result,
		final: req.final,
		ms: performance.now() - start,
		result: { signedJSON: resp.result, hash: resp.hash },
	};
}

// @ts-ignore
function handleRequest(req: any): any {
	switch (req.type as AggRequest) {
		case AggRequest.Ready:
			return ready(req);
		case AggRequest.Initialize:
			return respond(
				req,
				callGo(() => initialize(req)),
			);
		case AggRequest.Add:
			return respond(
				req,
				callGo(() => add(req)),
			);
		case AggRequest.Flush:
			return respond(
				req,
				callGo(() => doFlush(req)),
			);
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
