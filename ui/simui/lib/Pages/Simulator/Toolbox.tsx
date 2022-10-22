import { Classes, Button, Menu, MenuDivider, MenuItem, AnchorButton } from "@blueprintjs/core";
import { Popover2 } from "@blueprintjs/popover2";
import React from "react";
import { useTranslation } from "react-i18next";
import { useLocation } from "wouter";
import { RootState, useAppDispatch, useAppSelector } from "../../Stores/store";
import { userActions } from "~/Stores/userSlice";
import { SimWorkerOptions, ImportFromGOODDialog, ImportFromEnkaDialog } from "./Components";
import { runSim } from "~/Stores/viewerSlice";

type Props = {
  cfg: string;
  canRun?: boolean;
};

export const Toolbox = ({ cfg, canRun = true }: Props) => {
  const { t } = useTranslation();
  const [, setLocation] = useLocation();

  const [openImport, setOpenGOODImport] = React.useState<boolean>(false);
  const [openImportFromEnka, setOpenImportFromEnka] = React.useState<boolean>(false);
  const [openWorkers, setOpenWorkers] = React.useState<boolean>(false);
  const { ready, workers, settings } = useAppSelector((state: RootState) => {
    return {
      ready: state.app.ready,
      workers: state.app.workers,
      settings: state.user.settings ?? { showTips: false, showBuilder: false },
    };
  });

  const dispatch = useAppDispatch();
  const toggleTips = () => {
    dispatch(
      userActions.setUserSettings({
        showTips: !settings.showTips,
        showBuilder: settings.showBuilder,
      })
    );
  };

  const run = () => {
    dispatch(runSim(cfg));
    setLocation("/viewer/web");
  };

  const toggleBuilder = () => {
    dispatch(
      userActions.setUserSettings({
        showTips: settings.showTips,
        showBuilder: !settings.showBuilder,
      })
    );
  };

  const ToolMenu = (
    <Menu>
      <MenuItem
        icon="help"
        text={settings.showTips ? "Hide Tooltips" : "Show Tooltips"}
        onClick={toggleTips}
      />
      <MenuItem
        icon="people"
        text={settings.showBuilder ? "Hide Builder" : "Show Builder"}
        onClick={toggleBuilder}
      />
      <MenuDivider />
      <MenuItem icon="cut" text="Substat Snippets" disabled />

      <MenuItem text="Import from GO" icon="import" onClick={() => setOpenGOODImport(true)} />
      <MenuItem text="Import from Enka" icon="import" onClick={() => setOpenImportFromEnka(true)} />
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
            disabled={ready < workers || !canRun}
          >
            {ready < workers ? t("simple.loading_workers") : t("simple.run")}
          </Button>
        </div>
      </div>
      <ImportFromGOODDialog isOpen={openImport} onClose={() => setOpenGOODImport(false)} />
      <ImportFromEnkaDialog
        isOpen={openImportFromEnka}
        onClose={() => setOpenImportFromEnka(false)}
      />
      <SimWorkerOptions isOpen={openWorkers} onClose={() => setOpenWorkers(false)} />
    </div>
  );
};
