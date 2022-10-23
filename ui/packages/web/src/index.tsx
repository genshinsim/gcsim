import React from "react";
import ReactDOM from "react-dom/client";
import { App, SetExecutor } from "@gcsim/ui";
import { WasmExecutor } from "@gcsim/executors";

SetExecutor(new WasmExecutor());

ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);
