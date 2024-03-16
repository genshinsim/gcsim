import { Switch } from "@blueprintjs/core";
import React from "react";
import ServerMode from "./ServerMode";
import WasmMode from "./WasmMode";
import { useTranslation } from "react-i18next";

const serverModeKey = "use-server-mode";

const App = ({}) => {
  const { t } = useTranslation();
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
          {t<string>(serverMode ? "simple.server_mode_disable" : "simple.server_mode_enable")}
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
