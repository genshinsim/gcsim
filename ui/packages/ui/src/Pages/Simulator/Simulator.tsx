import { Callout, Spinner } from "@blueprintjs/core";
import React, { useEffect } from "react";
import { Viewport, SectionDivider } from "../../Components";
import { ActionList } from "./Components";
import { Team } from "./Team";
import { Trans, useTranslation } from "react-i18next";
import { Toolbox } from "./Toolbox";
import { ActionListTooltip, TeamBuilderTooltip } from "./Tooltips";
import { setTotalWorkers, updateCfg } from "../../Stores/appSlice";
import { useAppSelector, RootState, useAppDispatch } from "../../Stores/store";
import { Executor } from "@gcsim/executors";

export function Simulator({ pool }: { pool: Executor }) {
  const { t } = useTranslation();
  const dispatch = useAppDispatch();

  const { settings, ready, workers, cfg, cfgErr } = useAppSelector(
    (state: RootState) => {
      return {
        cfg: state.app.cfg,
        cfgErr: state.app.cfg_err,
        ready: state.app.ready,
        workers: state.app.workers,
        settings: state.user.settings ?? {
          showTips: false,
          showBuilder: false,
        },
      };
    }
  );
  useEffect(() => {
    dispatch(setTotalWorkers(pool, workers));
  }, []);
  useEffect(() => {
    //TODO: for whatever reason this is being called every single time a worker gets loaded
    //when it should only happen once?
    if (ready) {
      console.log("rerunning config on ready!");
      dispatch(updateCfg(cfg));
    }
  }, [ready]);

  // check worker ready state every 250ms so run button becomes available when workers do
  const [isReady, setReady] = React.useState<boolean>(false);
  useEffect(() => {
    const interval = setInterval(() => {
      setReady(pool.ready());
    }, 250);
    return () => clearInterval(interval);
  }, [pool]);

  if (ready === 0) {
    return (
      <Viewport>
        <Callout intent="primary" title={t("sim.loading_simulator_please")}>
          <Spinner />
        </Callout>
      </Viewport>
    );
  }

  return (
    <Viewport className="flex flex-col gap-2">
      <div className="flex flex-col gap-2">
        <div className="flex flex-col">
          {settings.showBuilder ? (
            <>
              <SectionDivider>
                <Trans>simple.team</Trans>
              </SectionDivider>
              <TeamBuilderTooltip />
              <Team />
            </>
          ) : null}

          <SectionDivider>
            <Trans>simple.action_list</Trans>
          </SectionDivider>

          <ActionListTooltip />

          <ActionList
            cfg={cfg}
            onChange={(v) => dispatch(updateCfg(v, false))}
          />

          <div className="sticky bottom-0 bg-bp-bg flex flex-col gap-y-1">
            {cfgErr !== "" ? (
              <div className="pl-2 pr-2 pt-2 mt-1">
                <Callout intent="warning" title="Error parsing config">
                  <pre className=" whitespace-pre-wrap">{cfgErr}</pre>
                </Callout>
              </div>
            ) : null}
            <Toolbox pool={pool} cfg={cfg} canRun={cfgErr === "" && isReady} />
          </div>
        </div>
      </div>
    </Viewport>
  );
}
