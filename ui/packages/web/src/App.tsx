import { Executor, WasmExecutor } from "@gcsim/executors";
import { UI } from "@gcsim/ui";

let exec: Executor | undefined;

const App = ({}) => {
  const supplier = () => {
    if (exec == null) {
      exec = new WasmExecutor();
    }
    return exec;
  };

  return <UI exec={supplier} />;
};

export default App;