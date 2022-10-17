import { Alert, Intent, Position, Toaster } from "@blueprintjs/core";
import { useState, useRef, useEffect } from "react";
import Config from "./Config";
import useLoadingToast from "./Components/LoadingToast";
import Results from "./Results";
import ViewerNav from "./Components/ViewerNav";
import { useLocation } from "wouter";
import { SimResults } from "./SimResults";
import { ViewTypes } from ".";

type ViewerProps = {
  data: SimResults | null;
  error: string | null;
  type: ViewTypes;
  redirect: string;
  cancel?: () => void;
  retry?: () => void;
};

export default ({ data, error, type, redirect, cancel, retry }: ViewerProps) => {
  const [tabId, setTabId] = useState("results");
  const loadingToast = useLoadingToast(
      type, error, cancel, data?.statistics?.iterations, data?.max_iterations);
  const copyToast = useRef<Toaster>(null);

  const tabs: { [k: string]: React.ReactNode } = {
    results: <Results data={data} />,
    config: <Config cfg={data?.config_file} />,
    analyze: <div></div>,
    debug: <div></div>,
  };

  // If we navigate away from the page, stop the execution
  useEffect(() => {
    return () => { cancel != null && cancel(); };
  }, [cancel]);

  // TODO: handle cases where schema is not compatible
  // - major version mismatch (dialog to rerun sim or reroute to legacy for V3)
  // - minor version mismatch (warning w/ option to resim or load page anyway)

  return (
    <div className="w-full bg-bp4-dark-gray-100">
      <div className="px-2 py-4 w-full 2xl:mx-auto 2xl:container">
        <ViewerNav tabState={[tabId, setTabId]} config={data?.config_file} copyToast={copyToast} />
      </div>
      <div className="pt-0 mt-0">
        {tabs[tabId]}
      </div>
      <Toaster ref={loadingToast} position={Position.TOP} />
      <Toaster ref={copyToast} position={Position.TOP_RIGHT} />
      <ErrorAlert msg={error} redirect={redirect} retry={retry} />
    </div>
  );
};

// TODO: Retry or Close buttons. Retry runs a callback?
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