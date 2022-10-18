import { Alert, Intent } from "@blueprintjs/core";
import { useState, useEffect } from "react";
import Config from "./Tabs/Config";
import Results from "./Tabs/Results";
import ViewerNav from "./Components/ViewerNav";
import { useLocation } from "wouter";
import { SimResults } from "./SimResults";
import { ViewTypes } from ".";
import LoadingToast from "./Components/LoadingToast";
import Debug, { useDebugParser, useDebugSettings } from "./Tabs/Debug";

type ViewerProps = {
  data: SimResults | null;
  error: string | null;
  type: ViewTypes;
  redirect: string;
  cancel?: () => void;
  retry?: () => void;
};

export default ({ data, error, type, redirect, cancel, retry }: ViewerProps) => {
  const [debugSettings, setDebugSettings] = useDebugSettings();
  const parsedDebug = useDebugParser(data, debugSettings);

  const [tabId, setTabId] = useState("results");
  const tabs: { [k: string]: React.ReactNode } = {
    results: <Results data={data} />,
    config: <Config cfg={data?.config_file} />,
    analyze: <div></div>,
    debug: <Debug data={data} parsed={parsedDebug} settingsState={[debugSettings, setDebugSettings]} />,
  };

  // If we navigate away from the page, stop the execution
  // TODO: push up to index?
  useEffect(() => {
    return () => { cancel != null && cancel(); };
  }, [cancel]);

  // TODO: handle cases where schema is not compatible
  // - major version mismatch (dialog to rerun sim or reroute to legacy for V3)
  // - minor version mismatch (warning w/ option to resim or load page anyway)

  return (
    <div className="flex flex-col w-full h-full bg-bp4-dark-gray-100">
      <div className="px-2 py-4 w-full 2xl:mx-auto 2xl:container">
        <ViewerNav tabState={[tabId, setTabId]} config={data?.config_file} />
      </div>
      <div className="basis-full pt-0 mt-0">
        {tabs[tabId]}
      </div>
      <LoadingToast
          type={type}
          error={error}
          current={data?.statistics?.iterations}
          total={data?.max_iterations}
          cancel={cancel} />
      <ErrorAlert msg={error} redirect={redirect} retry={retry} />
    </div>
  );
};

const ErrorAlert = ({ msg, redirect, retry }: {
    msg: string | null,
    redirect: string,
    retry?: () => void }) => {
  const [, setLocation] = useLocation();

  let cancelButtonText = undefined;
  let onCancel = undefined;
  if (retry != null) {
    cancelButtonText = "Retry";
    onCancel = () => retry();
  }

  return (
    <Alert
        isOpen={msg != null}
        onConfirm={() => setLocation(redirect)}
        onCancel={onCancel}
        canEscapeKeyCancel={false}
        canOutsideClickCancel={false}
        confirmButtonText="Close"
        cancelButtonText={cancelButtonText}
        intent={Intent.DANGER}>
      <p>{msg}</p>
    </Alert>
  );
};