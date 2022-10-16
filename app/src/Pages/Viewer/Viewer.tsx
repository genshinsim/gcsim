import { Alert, Intent, Position, Toaster } from "@blueprintjs/core";
import React, { RefObject } from "react";
import Config from "./Config";
import useLoadingToast from "./Components/LoadingToast";
import Results from "./Results";
import ViewerNav from "./Components/ViewerNav";
import { useLocation } from "wouter";
import { SimResults } from "./SimResults";

type ViewerProps = {
  data: SimResults | null;
  error: string | null;
};

export default ({ data, error }: ViewerProps) => {
  const loadingToast = useLoadingToast(428, 1000);
  const [tabId, setTabId] = React.useState("results");

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
      <div className="px-2 py-4 w-full 2xl:mx-auto 2xl:container">
        <ViewerNav tabState={[tabId, setTabId]} config={data?.config_file} />
      </div>
      <div className="pt-0 mt-0">
        {tabs[tabId]}
      </div>
      <Toaster ref={loadingToast} position={Position.TOP} maxToasts={1} />
      <ErrorAlert msg={error} loadingToast={loadingToast} />
    </div>
  );
};

// TODO: Retry or Close buttons. Retry runs a callback?
const ErrorAlert = ({ msg, loadingToast }:
      { msg: string | null, loadingToast: RefObject<Toaster> }) => {
  const [, setLocation] = useLocation();

  return (
    <Alert
        isOpen={msg != null}
        onOpening={() => loadingToast.current?.clear()}
        onConfirm={() => setLocation("/viewer")}
        canEscapeKeyCancel={false}
        canOutsideClickCancel={false}
        confirmButtonText="Close"
        intent={Intent.NONE}>
      <p>{msg}</p>
    </Alert>
  );
};