import React from "react";
import ReactDOM from "react-dom";
import { UI, SetExecutor } from "simui";
import { WorkerPool } from "./WorkerPool";
import "../node_modules/simui/dist/style.css";

let pool: WorkerPool = new WorkerPool();

SetExecutor(pool);

ReactDOM.render(
  <React.StrictMode>
    <UI />
  </React.StrictMode>,
  document.getElementById("root")
);
