import { Alert, Intent } from "@blueprintjs/core";
import { useEffect, useState } from "react";
import Config from "./Tabs/Config";
import Results from "./Tabs/Results";
import ViewerNav from "./Components/ViewerNav";
import { useLocation } from "wouter";
import { ResultSource } from ".";
import LoadingToast from "./Components/LoadingToast";
import Debug, { useDebug } from "./Tabs/Debug";
import { pool } from "~/Executor";
import { SimResults } from "~/Types";

type ViewerProps = {
  data: SimResults | null;
  error: string | null;
  src: ResultSource;
  redirect: string;
  retry?: () => void;
};

// The viewer is read-only against the data. Any mutations to the data (resim) must be performed
// above viewer in the hierarchy tree. The viewer can perform whatever additional calculations it
// wants (linreg, stat optimizations, etc) but these computations are *never* stored in the data and
// only exist as long as the page is loaded.
//
// The debug view is a partial "exception" to this rule. User can regenerate debug view on UI.
// This does not mutate the original data, but is stored as a new debug variable. The generated
// debug is only merged into the data when generating a share link (share link will have either
// no debug data or debug data from last generation).
export default ({ data, error, src, redirect, retry }: ViewerProps) => {
  const isRunning = useRunningState();
  const debug = useDebug(isRunning, data);

  const [tabId, setTabId] = useState("results");
  const tabs: { [k: string]: React.ReactNode } = {
    results: <Results data={data} />,
    config: <Config cfg={data?.config_file} />,
    analyze: <div></div>,
    debug: <Debug data={data} debug={debug} running={isRunning} />,
  };

  return (
    <div className="flex flex-col flex-grow w-full bg-bp4-dark-gray-100 pb-4">
      <div className="px-2 py-4 w-full 2xl:mx-auto 2xl:container">
        <ViewerNav
          tabState={[tabId, setTabId]}
          data={data}
          debug={debug.logs}
          running={isRunning}
        />
      </div>
      <div className="basis-full pt-0 mt-0">{tabs[tabId]}</div>
      <LoadingToast
        running={isRunning}
        src={src}
        error={error}
        current={data?.statistics?.iterations}
        total={data?.max_iterations}
      />
      <ErrorAlert msg={error} redirect={redirect} retry={retry} />
    </div>
  );
};

const ErrorAlert = ({
  msg,
  redirect,
  retry,
}: {
  msg: string | null;
  redirect: string;
  retry?: () => void;
}) => {
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
      intent={Intent.DANGER}
    >
      <p>{msg}</p>
    </Alert>
  );
};

function useRunningState(): boolean {
  const [isRunning, setRunning] = useState(true);

  useEffect(() => {
    const check = setInterval(() => {
      setRunning(pool.running());
    }, 250);
    return () => clearInterval(check);
  }, []);

  return isRunning;
}
