import { Alert, Intent } from "@blueprintjs/core";
import { useCallback, useState } from "react";
import ConfigUI, { useConfig } from "./Tabs/Config";
import SampleUI, { useSample } from "./Tabs/Sample";
import Results from "./Tabs/Results";
import ViewerNav from "./Components/ViewerNav";
import { useLocation } from "wouter";
import { ResultSource } from ".";
import LoadingToast from "./Components/LoadingToast";
import { SimResults } from "@gcsim/types";
import Warnings from "./Components/Warnings";
import { Executor, ExecutorSupplier } from "@gcsim/executors";

type ViewerProps = {
  running: boolean;
  data: SimResults | null;
  error: string | null;
  src: ResultSource;
  redirect: string;
  exec: ExecutorSupplier<Executor>;
  retry?: () => void;
};

// The viewer is read-only against the data. Any mutations to the data (resim) must be performed
// above viewer in the hierarchy tree. The viewer can perform whatever additional calculations it
// wants (linreg, stat optimizations, etc) but these computations are *never* stored in the data and
// only exist as long as the page is loaded.
export default ({ running, data, error, src, redirect, exec, retry }: ViewerProps) => {
  const sample = useSample(running, data);
  const config = useConfig(data, exec);

  const cancel = useCallback(() => exec().cancel(), [exec]);
  const sampler = useCallback((cfg: string, seed: string) => exec().sample(cfg, seed), [exec]);

  const [tabId, setTabId] = useState("results");
  const tabs: { [k: string]: React.ReactNode } = {
    results: <Results data={data} />,
    config: <ConfigUI config={config} running={running} />,
    analyze: <div></div>,
    sample: <SampleUI sampler={sampler} data={data} sample={sample} running={running} />,
  };

  return (
    <div className="flex flex-col flex-grow w-full bg-bp4-dark-gray-100 pb-6">
      <Warnings data={data} />
      <div className="px-6 py-4 w-full 2xl:mx-auto 2xl:container">
        <ViewerNav
          tabState={[tabId, setTabId]}
          data={data}
          running={running}
        />
      </div>
      <div className="basis-full pt-0 mt-0">{tabs[tabId]}</div>
      <LoadingToast
        cancel={cancel}
        running={running}
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

  let cancelButtonText: string | undefined;
  let onCancel: (() => void) | undefined;
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
