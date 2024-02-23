import { Switch } from "@blueprintjs/core";
import React from "react";
import ServerMode from "./ServerMode";
import WasmMode from "./WasmMode";

const serverModeKey = "use-server-mode";

const App = ({}) => {
  const [serverMode, setServerMode] = React.useState<boolean>((): boolean => {
    return localStorage.getItem(serverModeKey) === "true";
  });
  React.useEffect(() => {
    localStorage.setItem(serverModeKey, serverMode.toString());
  }, [serverMode]);

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
