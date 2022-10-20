import { Alert, Intent } from "@blueprintjs/core";
import { useState } from "react";
import Config from "./Tabs/Config";
import Results from "./Tabs/Results";
import ViewerNav from "./Components/ViewerNav";
import { useLocation } from "wouter";
import { SimResults } from "./SimResults";
import { ResultSource } from ".";
import LoadingToast from "./Components/LoadingToast";
import Debug, { useDebugParser, useDebugSettings } from "./Tabs/Debug";

type ViewerProps = {
  data: SimResults | null;
  error: string | null;
  src: ResultSource;
  redirect: string;
  retry?: () => void;
};

export default ({ data, error, src, redirect, retry }: ViewerProps) => {
  const [debugSettings, setDebugSettings] = useDebugSettings();
  const parsedDebug = useDebugParser(data, debugSettings);

  const [tabId, setTabId] = useState("results");
  const tabs: { [k: string]: React.ReactNode } = {
    results: <Results data={data} />,
    config: <Config cfg={data?.config_file} />,
    analyze: <div></div>,
    debug: (
      <Debug data={data} parsed={parsedDebug} settingsState={[debugSettings, setDebugSettings]} />
    ),
  };

  return (
    <div className="flex flex-col w-full h-full bg-bp4-dark-gray-100">
      <div className="px-2 py-4 w-full 2xl:mx-auto 2xl:container">
        <ViewerNav tabState={[tabId, setTabId]} data={data} />
      </div>
      <div className="basis-full pt-0 mt-0">
        {tabs[tabId]}
      </div>
      <LoadingToast
          src={src}
          error={error}
          current={data?.statistics?.iterations}
          total={data?.max_iterations} />
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