import { Alert, Callout, Intent, Position, Toaster } from "@blueprintjs/core";
import { useCallback, useMemo, useRef, useState } from "react";
import ConfigUI, { useConfig } from "./Tabs/Config";
import SampleUI, { useSample } from "./Tabs/Sample";
import Results from "./Tabs/Results";
import ViewerNav from "./Components/ViewerNav";
import { ResultSource } from ".";
import LoadingToast from "./Components/LoadingToast";
import { SimResults } from "@gcsim/types";
import Warnings from "./Components/Warnings";
import { Executor, ExecutorSupplier } from "@gcsim/executors";
import queryString from "query-string";
import { useHistory } from "react-router";
import { RootState, useAppSelector } from "@ui/Stores/store";
import CopyToClipboard from "@ui/Components/Buttons/CopyToClipboard";
import SendToSimulator from "@ui/Components/Buttons/SendToSimulator";
import { useTranslation } from "react-i18next";

type ViewerProps = {
  running: boolean;
  data: SimResults | null;
  hash: string | null;
  recoveryConfig: string | null;
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
export default ({ running, data, hash = "", recoveryConfig, error, src, redirect, exec, retry }: ViewerProps) => {
  const { t } = useTranslation();
  const parsed = queryString.parse(location.hash);
  const [tabId, setTabId] = useState((parsed.tab as string) ?? "results");

  const cancel = useCallback(() => exec().cancel(), [exec]);
  const sampler = useCallback((cfg: string, seed: string) => exec().sample(cfg, seed), [exec]);
  const resetTab = useCallback(() => setTabId("results"), []);

  const { sampleOnLoad } = useAppSelector((state: RootState) => {
    return {
      sampleOnLoad: state.app.sampleOnLoad,
    };
  });

  const sample = useSample(running, data, sampleOnLoad, sampler);
  const config = useConfig(data, exec);
  const names = useMemo(
      () => data?.character_details?.map(c => t<string>("character_names." + c.name, { ns: "game" })), [data?.character_details, t]);

  const tabs: { [k: string]: React.ReactNode } = {
    results: <Results data={data} running={running} names={names} />,
    config: <ConfigUI config={config} running={running} resetTab={resetTab} />,
    analyze: <div></div>,
    sample: <SampleUI sampler={sampler} data={data} sample={sample} running={running} />,
  };

  return (
    <div className="flex flex-col flex-grow w-full pb-6">
      <Warnings data={data} />
      <div className="px-2 py-4 w-full 2xl:mx-auto 2xl:container">
        <ViewerNav
          hash={hash}
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
        total={data?.simulator_settings?.iterations}
      />
      <ErrorAlert msg={error} recoveryConfig={recoveryConfig} redirect={redirect} retry={retry} />
    </div>
  );
};

const ErrorAlert = ({
      msg,
      recoveryConfig,
      redirect,
      retry,
    }: {
      msg: string | null;
      recoveryConfig: string | null;
      redirect: string;
      retry?: () => void;
    }) => {
  const { t } = useTranslation();
  const copyToast = useRef<Toaster>(null);
  const history = useHistory();

  let cancelButtonText: string | undefined;
  let onCancel: (() => void) | undefined;
  if (retry != null) {
    cancelButtonText = t<string>("viewer.retry");
    onCancel = () => retry();
  }

  return (
    <Alert
      isOpen={msg != null}
      onConfirm={() => history.push(redirect)}
      onCancel={onCancel}
      canEscapeKeyCancel={false}
      canOutsideClickCancel={false}
      confirmButtonText={t<string>("viewer.return_to_sim")}
      cancelButtonText={cancelButtonText}
      intent={Intent.DANGER}
    >
      <div className="flex flex-col gap-2 mb-1">
        <Callout intent={Intent.DANGER} title={t<string>("viewer.error_encountered")}>
          <pre className="whitespace-pre-wrap pl-5">{msg}</pre>
        </Callout>
        {recoveryConfig != null ? (
          <>
            <CopyToClipboard
              copyToast={copyToast}
              config={recoveryConfig}
              className="hidden ml-[7px] sm:flex"
            />
            <SendToSimulator config={recoveryConfig} />
          </>
        ) : null}
      </div>
      <Toaster ref={copyToast} position={Position.TOP_RIGHT} />
    </Alert>
  );
};
