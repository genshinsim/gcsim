import { Alert, Colors, Intent } from "@blueprintjs/core";
import React from "react";
import Config from "./Config";
import useLoadingToast, { LoadingToaster } from "./Components/LoadingToast";
import Results from "./Results";
import ViewerNav from "./Components/ViewerNav";
import { useLocation } from "wouter";

type ViewerProps = {
  data: any | null;
  error: any | null;
};

export default ({ data, error }: ViewerProps) => {
  const [tabId, setTabId] = React.useState("results");
  const isLoaded = useLoadingToast(data);

  const tabs: { [k: string]: React.ReactNode } = {
    results: <Results data={data} />,
    config: <Config cfg={data?.config_file} />,
    analyze: <div></div>,
    debug: <div></div>,
  };

  // TODO: handle cases where schema is not compatible
  // - major version mismatch (dialog to rerun sim or reroute to legacy for V3)
  // - minor version mismatch (warning w/ option to resim or load page anyway)

  return (
    <div className="w-full bg-bp4-dark-gray-100">
      <div className="p-4 w-full 2xl:mx-auto 2xl:container">
        <ViewerNav isLoaded={isLoaded} tabState={[tabId, setTabId]} config={data?.config_file} />
      </div>
      <div className="pt-0 mt-0">
        {tabs[tabId]}
      </div>
      <ErrorAlert msg={error} />
    </div>
  );
};

const ErrorAlert = ({ msg }: { msg: string | null }) => {
  const [_, setLocation] = useLocation();

  return (
    <Alert
        isOpen={msg != null}
        onOpening={() => LoadingToaster.clear()}
        onConfirm={() => setLocation("/viewer")}
        canEscapeKeyCancel={false}
        canOutsideClickCancel={false}
        confirmButtonText="Okay"
        intent={Intent.DANGER}>
      <p>{msg}</p>
    </Alert>
  );
};