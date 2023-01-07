import { FormGroup, NumericInput } from "@blueprintjs/core";
import { ExecutorSupplier, WasmExecutor } from "@gcsim/executors";
import { UI } from "@gcsim/ui";
import { useLocalStorage } from "@gcsim/utils";
import { useRef } from "react";

const minWorkers = 1;
const maxWorkers = 30;

let exec: WasmExecutor | undefined;

function wasmLocation() {
  if (import.meta.env.PROD) {
    return "/api/wasm/"
        + import.meta.env.VITE_GIT_BRANCH + "/"
        + import.meta.env.VITE_GIT_COMMIT_HASH + "/"
        + "main.wasm";
  }
  return "/main.wasm";
}

const App = ({}) => {
  const [workers, setWorkers] = useLocalStorage<number>("wasm-num-workers", 3);

  const supplier = useRef<ExecutorSupplier<WasmExecutor>>(() => {
    if (exec == null) {
      exec = new WasmExecutor(wasmLocation());
      exec.setWorkerCount(workers);
    }
    return exec;
  });

  const updateWorkers = (num: number) => {
    num = Math.min(Math.max(num, minWorkers), maxWorkers);
    setWorkers(num);
    supplier.current().setWorkerCount(num);
  };

  return (
    <UI
        exec={supplier.current}
        gitCommit={import.meta.env.VITE_GIT_COMMIT_HASH}
        mode={import.meta.env.MODE}>
      <FormGroup className="!m-0" label="Workers">
        <NumericInput
          value={workers}
          onValueChange={(v) => updateWorkers(v)}
          min={minWorkers}
          max={maxWorkers}
          fill={true}
        />
      </FormGroup>
    </UI>
  );
};

export default App;