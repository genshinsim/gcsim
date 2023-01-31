import React, { useEffect, useMemo, useState } from 'react';
import Editor from "react-simple-code-editor";
import { Button, Callout, Intent, NonIdealState, Spinner, SpinnerSize } from "@blueprintjs/core";
import { SimResults } from '@gcsim/types';
import ExecutorSettingsButton from '../../../Components/Buttons/ExecutorSettingsButton';
import { Executor, ExecutorSupplier } from '@gcsim/executors';
import { runSim } from '../../Simulator/Toolbox';
import { useAppDispatch } from '../../../Stores/store';
import { useConfigValidateListener } from '../../Simulator';

//@ts-ignore
import { highlight, languages } from "prismjs/components/prism-core";
import "prismjs/components/prism-gcsim";
import "prismjs/themes/prism-tomorrow.css";
import { useHistory } from 'react-router';

type UseConfigData = {
  cfg?: string;
  error: string;
  isReady: boolean | null;
  validated: boolean;
  modified: boolean;
  exec: ExecutorSupplier<Executor>;
  setCfg: (cfg: string) => void;
}

type ConfigProps = {
  config: UseConfigData;
  running: boolean;
  resetTab: () => void;
};

const ConfigUI = ({ config, running, resetTab }: ConfigProps) => {
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
              text="Rerun"
              intent={Intent.SUCCESS}
              disabled={config.error !== "" || (!config.validated && config.modified)}
              loading={!config.isReady || running}
              className="basis-1/2"
              onClick={() => {
                dispatch(runSim(config.exec(), config.cfg ?? ""));
                resetTab();
                history.push("/viewer/web");
              }} />
        </div>
        <Error error={config.error} cfg={config.cfg} />
      </div>
      <div>
        <Editor
            value={config.cfg}
            onValueChange={(c) => config.setCfg(c)}
            textareaId="codeArea"
            className="editor"
            highlight={(code) =>
              highlight(code, languages.gcsim)
                .split("\n")
                .map(
                  (line: string, i: number) =>
                    `<span class='editorLineNumber'>${i + 1}</span>${line}`
                )
                .join("\n")
            }
            insertSpaces
            padding={10}
            style={{
              fontFamily: '"Fira code", "Fira Mono", monospace',
              fontSize: 14,
              backgroundColor: "rgb(45 45 45)",
            }} />
      </div>
    </div>
  );
};

const Error = ({ error, cfg }: { error: string, cfg: string}) => {
  if (error === "" || cfg === "") {
    return null;
  }
  return (
    <div className="px-6 pt-4">
      <Callout intent={Intent.DANGER} title="Error: Config Invalid">
        <pre className="whitespace-pre-wrap pl-5">{error}</pre>
      </Callout>
    </div>
  );
};

export function useConfig(data: SimResults | null, exec: ExecutorSupplier<Executor>): UseConfigData {
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
      setReady(exec().ready());
    }, 250);
    return () => clearInterval(interval);
  }, [exec]);

  // will detect changes in the redux config and validate with the executor
  // validated == true means we had a successful validation check run, not that it is valid
  const validated = useConfigValidateListener(exec, cfg ?? "", modified, setErr);

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