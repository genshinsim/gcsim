import { Classes, Button, Menu, MenuDivider, MenuItem } from "@blueprintjs/core";
import { Popover2 } from "@blueprintjs/popover2";
import React from "react";
import { useTranslation } from "react-i18next";
import { useLocation } from "wouter";
import { AppThunk, RootState, useAppDispatch, useAppSelector } from "../../Stores/store";
import { userActions } from "../../Stores/userSlice";
import { ImportFromGOODDialog, ImportFromEnkaDialog } from "./Components";
import { viewerActions } from "../../Stores/viewerSlice";
import { Executor, ExecutorSupplier } from "@gcsim/executors";
import ExecutorSettingsButton from "../../Components/ExecutorSettingsButton";
import { throttle } from "lodash-es";
import { SimResults } from "@gcsim/types";
import { VIEWER_THROTTLE } from "../Viewer";

type Props = {
  exec: ExecutorSupplier;
  cfg: string;
  canRun?: boolean;
};

function runSim(pool: Executor, cfg: string): AppThunk {
  return function (dispatch) {
    console.log("starting run");
    dispatch(viewerActions.start());

    const updateResult = throttle(
      (res: SimResults) => {
        dispatch(viewerActions.setResult({ data: res }));
      },
      VIEWER_THROTTLE,
      { leading: true, trailing: true }
    );

    pool.run(cfg, (result) => {
      updateResult(result);
    }).catch((err) => {
      dispatch(viewerActions.setError({ error: err }));
    });
  };
}

export const Toolbox = ({ exec, cfg, canRun = true }: Props) => {
  const { t } = useTranslation();
  const [, setLocation] = useLocation();

  const [openImport, setOpenGOODImport] = React.useState<boolean>(false);
  const [openImportFromEnka, setOpenImportFromEnka] = React.useState<boolean>(false);
  const { settings } = useAppSelector((state: RootState) => {
    return {
      settings: state.user.settings ?? { showTips: false, showBuilder: false },
    };
  });

  const dispatch = useAppDispatch();
  const toggleTips = () => {
    dispatch(
      userActions.setUserSettings({
        showTips: settings.showTips,
        showBuilder: settings.showBuilder,
      })
    );
  };

  const run = () => {
    dispatch(runSim(exec(), cfg));
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
        <ExecutorSettingsButton />
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
              loading={!canRun}
              text={t<string>("simple.run")} />
        </div>
      </div>
      <ImportFromGOODDialog isOpen={openImport} onClose={() => setOpenGOODImport(false)} />
      <ImportFromEnkaDialog
        isOpen={openImportFromEnka}
        onClose={() => setOpenImportFromEnka(false)}
      />
    </div>
  );
};
