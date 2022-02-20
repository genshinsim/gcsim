// let url = process.env.PUBLIC_URL + "/wasm_exec.js";
// @ts-ignore
self.importScripts("wasm_exec.js");

if (!WebAssembly.instantiateStreaming) {
  // polyfill
  WebAssembly.instantiateStreaming = async (resp, importObject) => {
    const source = await (await resp).arrayBuffer();
    return await WebAssembly.instantiate(source, importObject);
  };
}

//@ts-ignore
const go = new Go(); // Defined in wasm_exec.js. Don't forget to add this in your index.html.

declare function sim(): string;
declare function debug(): string;
declare function collect(data: string): string;
declare function setcfg(content: string): string;
declare function version(): string;

let inst: WebAssembly.Instance;
WebAssembly.instantiateStreaming(fetch("/main.wasm"), go.importObject)
  .then((result) => {
    inst = result.instance;
    go.run(inst);
    console.log("worker loaded ok");
    postMessage("ready");
  })
  .catch((err) => {
    console.error(err);
  });

onmessage = async (ev) => {
  // console.log(ev.data);
  switch (ev.data.cmd) {
    case "run":
      // console.log("running...");
      const simres = sim();
      postMessage(simres);
      break;
    case "collect":
      const res = collect(ev.data.payload);
      postMessage(res);
      break;
    case "debug":
      const d = debug();
      // console.log("debug run done");
      postMessage(d);
      break;
    case "version":
      const v = version();
      postMessage(v);
      break;
    case "cfg":
      const ok = setcfg(ev.data.payload);
      // console.log("configs set: " + ok);
      postMessage(ok);
      break;
  }
};
