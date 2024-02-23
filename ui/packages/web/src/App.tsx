import { Switch } from "@blueprintjs/core";
import React from "react";
import ServerMode from "./ServerMode";
import WasmMode from "./WasmMode";

const App = ({}) => {
  const [serverMode, setServerMode] = React.useState(true);

  const children = (
    <Switch
      checked={serverMode}
      onChange={() => setServerMode(!serverMode)}
      labelElement={
        <span>
          {`${serverMode ? "Turn off" : "Use"} server mode`} (advanced users
          only)
        </span>
      }
    />
  );

  return (
    <>
      {serverMode ? (
        <ServerMode children={children} />
      ) : (
        <WasmMode children={children} />
      )}
    </>
  );
};

export default App;
