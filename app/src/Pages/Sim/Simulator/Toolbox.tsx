import {
  Classes,
  Button,
  Menu,
  MenuDivider,
  MenuItem,
  AnchorButton,
} from "@blueprintjs/core";
import { Popover2 } from "@blueprintjs/popover2";
import React from "react";
import { useTranslation } from "react-i18next";
import { ImportFromGOODDialog } from "~src/Components/GOOD";
import { useAppSelector, RootState, useAppDispatch } from "~src/store";
import { SimWorkerOptions } from "../Components";
import { SimProgress } from "../Components/SimProgress";
import { runSim } from "../exec";
import { simActions } from "../simSlice";

export const Toolbox = ({ canRun = true }: { canRun?: boolean }) => {
  let { t } = useTranslation();

  const { ready, workers, cfg, run_stats, showTips, showBuilder } =
    useAppSelector((state: RootState) => {
      return {
        ready: state.sim.ready,
        workers: state.sim.workers,
        cfg: state.sim.cfg,
        run_stats: state.sim.run,
        showTips: state.sim.showTips,
        showBuilder: state.sim.showBuilder,
      };
    });
  const [openImport, setOpenImport] = React.useState<boolean>(false);
  const [openProgress, setOpenProgress] = React.useState<boolean>(false);
  const [openWorkers, setOpenWorkers] = React.useState<boolean>(false);

  const dispatch = useAppDispatch();

  const run = () => {
    dispatch(runSim(cfg));
    setOpenProgress(true);
  };

  const toggleTips = () => {
    dispatch(simActions.setShowTips(!showTips));
  };

  const toggleBuilder = () => {
    dispatch(simActions.setShowBuilder(!showBuilder));
  };

  const ToolMenu = (
    <Menu>
      <MenuItem
        icon="help"
        text={showTips ? "Hide Tooltips" : "Show Tooltips"}
        onClick={toggleTips}
      />
      <MenuItem
        icon="people"
        text={showBuilder ? "Hide Builder" : "Show Builder"}
        onClick={toggleBuilder}
      />
      <MenuDivider />
      <MenuItem icon="cut" text="Substat Snippets" disabled />

      <MenuItem
        text="Import from GO"
        icon="import"
        onClick={() => setOpenImport(true)}
      />
    </Menu>
  );

  return (
    <div className="p-2 wide:ml-2 wide:mr-2 flex flex-row flex-wrap place-items-center gap-x-1 gap-y-1">
      <div className="basis-full wide:basis-0 flex-grow p-1 flex flex-row items-center">
        <div className="pr-2">{`${t("simple.workers_available")}${ready}`}</div>
        <Button icon="edit" minimal onClick={() => setOpenWorkers(true)} />
      </div>
      <div className="basis-full wide:basis-2/3 p-1 flex flex-row flex-wrap">
        <Popover2
          content={ToolMenu}
          placement="top"
          className="basis-full md:basis-1/2"
          popoverClassName={Classes.POPOVER_DISMISS}
        >
          <Button icon="wrench" fill>
            Tools
          </Button>
        </Popover2>
        <div className="basis-full md:basis-1/2">
          <Button
            icon="play"
            fill
            intent="primary"
            onClick={run}
            disabled={ready < workers || run_stats.progress !== -1 || !canRun}
          >
            {ready < workers ? t("simple.loading_workers") : t("simple.run")}
          </Button>
        </div>
      </div>
      <SimProgress
        isOpen={openProgress}
        onClose={() => setOpenProgress(false)}
      />
      <ImportFromGOODDialog
        isOpen={openImport}
        onClose={() => setOpenImport(false)}
      />
      <SimWorkerOptions
        isOpen={openWorkers}
        onClose={() => setOpenWorkers(false)}
      />
    </div>
  );
};
