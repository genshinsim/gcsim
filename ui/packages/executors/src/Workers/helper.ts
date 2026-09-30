/* eslint-disable @typescript-eslint/no-explicit-any */
// @ts-ignore
self.importScripts("/wasm_exec.js");

let readyState = false;
let loadError: string | null = null;
// @ts-ignore
let go: Go;
let compiled: Promise<WebAssembly.Module> | null = null;

// @ts-ignore
function ready(req: { wasm: string }) {
	go = new Go();
	compiled = compileWasm(req.wasm);
	compiled
		.then((module) => {
			// the executor can now have the module shared with the aggregator and the sim workers
			postMessage({ type: HelpResponse.Ready });
			return WebAssembly.instantiate(module, go.importObject);
		})
		.then((instance) => {
			go.run(instance);
			console.log("helper loaded okay");
			readyState = true;
			processQueue();
		})
		.catch((e) => {
			console.error(e);
			loadError = e instanceof Error ? e.message : "Unknown Error";
			processQueue();
			postMessage({
				type: HelpResponse.Failed,
				reason: loadError,
				fatal: true,
			});
		});
}

// Handles a request that calls into Go. The Go functions return their errors, so if the call
// throws, or the Go program exits during it (a fatal error such as running out of memory), this
// instance is unusable: the fatal response makes the executor replace the helper.
// @ts-ignore
function callGo(req: { id: number }, handle: () => any): any {
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
		type: HelpResponse.Failed,
		reason: `The helper crashed: ${reason}`,
		fatal: true,
		id: req.id,
	};
}

// Sends the compiled module to an aggregator or a sim worker, over a MessageChannel the executor
// set up between the two, once that worker asks for it.
// - The module must not go through the page. Under the page's CSP (script-src without
//   'wasm-unsafe-eval'), Firefox refuses to deserialize a WebAssembly.Module there and the page
//   only gets a messageerror. Workers loaded from their own URL don't inherit the page's CSP.
// - Firefox also can't deliver a module to a port that is still being transferred: the
//   receiver gets a messageerror. The worker's request shows its end has arrived.
function share(req: { id: number; port: MessagePort }) {
	req.port.onmessage = () => {
		compiled?.then(
			(module) => {
				req.port.postMessage(module);
				postMessage({ type: HelpResponse.Shared, id: req.id });
			},
			// the failed load is fatal, and the executor fails this request with it
			() => {},
		);
	};
}

function compileWasm(url: string): Promise<WebAssembly.Module> {
	if (!WebAssembly.compileStreaming) {
		// polyfill
		return fetch(url)
			.then((resp) => resp.arrayBuffer())
			.then((source) => WebAssembly.compile(source));
	}
	return WebAssembly.compileStreaming(fetch(url));
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

// @ts-ignore
function handleRequest(req: any): any {
	switch (req.type as HelpRequest) {
		case HelpRequest.Validate:
			return callGo(req, () => validate(req));
		case HelpRequest.Sample:
			return callGo(req, () => doSample(req));
		default:
			console.error("helper - unknown request: ", req);
			throw new Error("helper unknown request");
	}
}

const queue: MessageEvent<any>[] = [];
self.onmessage = (ev) => {
	if (ev.data.type === HelpRequest.Ready) {
		ready(ev.data);
		return;
	}
	// needs the compiled module only, not the Go program
	if (ev.data.type === HelpRequest.Share) {
		share(ev.data);
		return;
	}

	queue.push(ev);
	processQueue();
};

// Requests that arrive before the wasm is loaded wait in the queue; ready() drains it.
function processQueue() {
	if (!readyState && loadError == null) {
		return;
	}

	for (let event = queue.shift(); event; event = queue.shift()) {
		if (loadError != null) {
			postMessage({
				type: HelpResponse.Failed,
				reason: loadError,
				id: event.data.id,
			});
			continue;
		}
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
	Share = "share",
}

enum HelpResponse {
	Failed = "failed",
	Ready = "ready",
	Validate = "validated",
	Sample = "sample",
	Shared = "shared",
}
