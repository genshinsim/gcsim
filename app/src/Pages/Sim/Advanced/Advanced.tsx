import { Button, Callout, Collapse } from "@blueprintjs/core";
import React from "react";
import { SectionDivider } from "~src/Components/SectionDivider";
import { Viewport } from "~src/Components/Viewport";
import { RootState, useAppDispatch, useAppSelector } from "~src/store";
import { simActions, updateAdvConfig } from "..";
import { ActionList, SimOptions } from "../Components";
import { SimProgress } from "../Components/SimProgress";
import { runSim } from "../exec";
import { Trans, useTranslation } from "react-i18next";

export function Advanced() {
  let { t } = useTranslation();

  const { ready, workers, cfg, cfg_err, runState } = useAppSelector(
    (state: RootState) => {
      return {
        ready: state.sim.ready,
        workers: state.sim.workers,
        cfg: state.sim.advanced_cfg,
        cfg_err: state.sim.adv_cfg_err,
        runState: state.sim.run,
      };
    }
  );
  const dispatch = useAppDispatch();
  const [open, setOpen] = React.useState<boolean>(false);
  const [showOptions, setShowOptions] = React.useState<boolean>(false);

  const run = () => {
    dispatch(runSim(cfg));
    setOpen(true);
  };

  return (
    <Viewport className="flex flex-col gap-2">
      <div className="flex flex-col">
        <SectionDivider>
          <Trans>advanced.action_list</Trans>
        </SectionDivider>
        <ActionList cfg={cfg} onChange={(v) => dispatch(updateAdvConfig(v))} />
        <SectionDivider>
          <Trans>advanced.helpers</Trans>
        </SectionDivider>
        <div className="p-2">
          <Button disabled>
            <Trans>advanced.substat_helper</Trans>
          </Button>
        </div>
        <SectionDivider>
          <Trans>advanced.sim_options</Trans>
        </SectionDivider>
        <div className="ml-auto mr-2">
          <Button icon="edit" onClick={() => setShowOptions(!showOptions)}>
            {showOptions ? t("advanced.hide") : t("advanced.show")}
          </Button>
        </div>
      </div>
      <div className="sticky bottom-0 bg-bp-bg p-2 wide:ml-2 wide:mr-2 flex flex-col ">
        {cfg_err !== "" ? (
          <div className="basis-full p-1">
            <Callout intent="warning" title="Error parsing config">
              <pre className=" whitespace-pre-wrap">{cfg_err}</pre>
            </Callout>
          </div>
        ) : null}

        <div className="flex flex-row  flex-wrap place-items-center gap-x-1 gap-y-1">
          <div className="basis-full wide:basis-0 flex-grow p-1">
            {`${t("advanced.workers_available")}${ready}`}
          </div>
          <div className="basis-full wide:basis-1/3 p-1">
            <Button
              icon="play"
              fill
              intent="primary"
              onClick={run}
              disabled={
                ready < workers || runState.progress !== -1 || cfg_err !== ""
              }
            >
              {ready < workers
                ? t("advanced.loading_workers")
                : t("advanced.run")}
            </Button>
          </div>
        </div>
      </div>
      <SimProgress isOpen={open} onClose={() => setOpen(false)} />
    </Viewport>
  );
}
