import React from "react";
import ReactDOM from "react-dom/client";
import { UI } from "@gcsim/ui";
import { WasmExecutor } from "@gcsim/executors";

const exec = new WasmExecutor();
const supplier = () => exec;

ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
  <React.StrictMode>
    <UI exec={supplier} />
  </React.StrictMode>
);