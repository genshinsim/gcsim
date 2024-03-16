import {
  Button,
  Callout,
  Intent,
  NonIdealState,
  Spinner,
  SpinnerSize,
} from "@blueprintjs/core";
import { Executor, ExecutorSupplier } from "@gcsim/executors";
import { SimResults } from "@gcsim/types";
import React, { useEffect, useMemo, useState } from "react";
import ExecutorSettingsButton from "../../../Components/Buttons/ExecutorSettingsButton";
import { useAppDispatch } from "../../../Stores/store";
import { useConfigValidateListener } from "../../Simulator";
import { runSim } from "../../Simulator/Toolbox";

import { ConfigEditor } from "@ui/Components";
import { useTranslation } from "react-i18next";
import { useHistory } from "react-router";

type UseConfigData = {
  cfg?: string;
  error: string;
  isReady: boolean | null;
  validated: boolean;
  modified: boolean;
  exec: ExecutorSupplier<Executor>;
  setCfg: (cfg: string) => void;
};

type ConfigProps = {
  config: UseConfigData;
  running: boolean;
  resetTab: () => void;
};

const ConfigUI = ({ config, running, resetTab }: ConfigProps) => {
  const { t } = useTranslation();
  const dispatch = useAppDispatch();
  const history = useHistory();

  if (config.cfg == null) {
    return <NonIdealState icon={<Spinner size={SpinnerSize.LARGE} />} />;
  }

  return (
    <div className="w-full 2xl:mx-auto 2xl:container -mt-4 px-2">
      <div className="sticky top-0 bg-bp4-dark-gray-100 py-4 z-10">
        <div className="flex gap-2 justify-center">
          <ExecutorSettingsButton />
          <Button
            icon="refresh"
            text={t<string>("viewer.rerun")}
            intent={Intent.SUCCESS}
            disabled={
              config.error !== "" || (!config.validated && config.modified)
            }
            loading={!config.isReady || running}
            className="basis-1/2"
            onClick={() => {
              dispatch(runSim(config.exec(), config.cfg ?? ""));
              resetTab();
              history.push("/web");
            }}
          />
        </div>
        <Error error={config.error} cfg={config.cfg} />
      </div>
      <div>
        <ConfigEditor cfg={config.cfg} onChange={config.setCfg} />
      </div>
    </div>
  );
};

const Error = ({ error, cfg }: { error: string; cfg: string }) => {
  const { t } = useTranslation();
  if (error === "" || cfg === "") {
    return null;
  }
  return (
    <div className="px-6 pt-4">
      <Callout
        intent={Intent.DANGER}
        title={
          t<string>("viewer.error_encountered") +
          +t<string>("viewer.config_invalid")
        }
      >
        <pre className="whitespace-pre-wrap pl-5">{error}</pre>
      </Callout>
    </div>
  );
};

export function useConfig(
  data: SimResults | null,
  exec: ExecutorSupplier<Executor>
): UseConfigData {
  const [cfg, setCfg] = useState(data?.config_file);
  const [isReady, setReady] = useState<boolean | null>(null);
  const [err, setErr] = useState("");
  const [modified, setModified] = useState<boolean>(false);

  const updateCfg = (newCfg: string) => {
    setCfg(newCfg);
    setModified(true);
  };

  // reset config file every time it changes from results
  useEffect(() => {
    setCfg(data?.config_file);
  }, [data?.config_file]);

  // check worker ready state every 250ms so run button becomes available when workers do
  useEffect(() => {
    const interval = setInterval(() => {
      exec()
        .ready()
        .then((res) => setReady(res));
    }, 250);
    return () => clearInterval(interval);
  }, [exec]);

  // will detect changes in the redux config and validate with the executor
  // validated == true means we had a successful validation check run, not that it is valid
  const validated = useConfigValidateListener(
    exec,
    cfg ?? "",
    modified,
    setErr
  );

  return useMemo(() => {
    return {
      cfg: cfg,
      error: err,
      isReady: isReady,
      validated: validated,
      modified: modified,
      exec: exec,
      setCfg: updateCfg,
    };
  }, [cfg, err, exec, isReady, modified, validated]);
}

export default React.memo(ConfigUI);
