import {
  Button,
  ButtonGroup,
  Callout,
  Checkbox,
  Classes,
  Dialog,
  Tab,
  Tabs,
} from "@blueprintjs/core";
import React from "react";
import { Config } from "./Config";
import { SimResults } from "./DataType";
import { Debugger } from "./DebugView";
import { Details } from "./Details";
import { Options, OptionsProp } from "./Options";
import { DebugRow, parseLog } from "./parse";
import Summary from "./Summary";
import Share, { AppToaster } from "./Share";
import { parseLogV2 } from "./parsev2";
import { Trans, useTranslation } from "react-i18next";
import { useLocation } from "wouter";
import { updateCfg } from "~src/Pages/Sim";
import { useAppDispatch } from "~src/store";

const opts = [
  "procs",
  "damage",
  "pre_damage_mods",
  "hurt",
  "heal",
  "calc",
  "reaction",
  "element",
  "snapshot",
  "snapshot_mods",
  "status",
  "action",
  "queue",
  "energy",
  "character",
  "enemy",
  "hook",
  "sim",
  "task",
  "artifact",
  "weapon",
  "shield",
  "construct",
  "icd",
];

const defOpts = [
  "damage",
  "element",
  "action",
  "energy",
  "pre_damage_mods",
  "status",
];

type ViewProps = {
  classes?: string;
  selected: string[];
  handleSetSelected: (next: string[]) => void;
  data: SimResults;
  parsed: DebugRow[];
  handleClose: () => void;
};

let msgs: { [key: number]: any } = {};
const extractMsgData = (data: DebugRow[]) => {
  data.map((row, i) => {
    let results: string[] = [];

    row.slots.map((slot, ci) => {
      slot.map((e, ei) => {
        results.push(e.msg);
      })
    })

    msgs[i] = results;
  })

  console.log(msgs);
}

const LOCALSTORAGE_KEY = "gcsim-viewer-cpy-cfg-settings";

function ViewOnly(props: ViewProps) {
  let { t } = useTranslation();
  const [open, setOpen] = React.useState<boolean>(false);
  const [tabID, setTabID] = React.useState<string>("result");
  const [optOpen, setOptOpen] = React.useState<boolean>(false);

  extractMsgData(props.parsed);

  const handleTabChange = (next: string) => {
    setTabID(next);
  };

  const optProps: OptionsProp = {
    isOpen: optOpen,
    handleClose: () => {
      setOptOpen(false);
    },
    handleToggle: (t: string) => {
      const i = props.selected.indexOf(t);
      let next = [...props.selected];
      if (i === -1) {
        next.push(t);
      } else {
        next.splice(i, 1);
      }
      props.handleSetSelected(next);
    },
    handleClear: () => {
      props.handleSetSelected([]);
    },
    handleResetDefault: () => {
      props.handleSetSelected(defOpts);
    },
    selected: props.selected,
    options: opts,
  };

  function copyToClipboard() {
    navigator.clipboard.writeText(props.data.config_file).then(() => {
      AppToaster.show({
        message: t("viewer.copied_to_clipboard"),
        intent: "success",
      });
    });
    // TODO: Need to add a blueprintjs Toaster for ephemeral confirmation box
  }

  return (
    <div
      className={props.classes + " p-4 rounded-lg bg-gray-800 flex flex-col"}
    >
      <div className="flex flex-row  bg-gray-800 ">
        <Tabs
          selectedTabId={tabID}
          onChange={handleTabChange}
          className="w-full"
        >
          <Tab
            id="result"
            title={t("viewer.summary")}
            className="focus:outline-none"
          />
          <Tab
            id="details"
            title={t("viewer.details")}
            className="focus:outline-none"
          />
          <Tab
            id="config"
            title={t("viewer.config")}
            className="focus:outline-none"
          />
          <Tab
            id="debug"
            title={t("viewer.debug")}
            className="focus:outline-none"
          />
          <Tab
            id="share"
            title={t("viewer.share")}
            className="focus:outline-none"
          />
          <Tabs.Expander />
          <ButtonGroup>
            <Button onClick={copyToClipboard} icon="clipboard">
              <Trans>viewer.copy</Trans>
            </Button>
            <Button onClick={() => setOpen(true)} icon="send-to">
              <Trans>viewer.send_to_simulator</Trans>
            </Button>
          </ButtonGroup>

          <Button icon="cross" intent="danger" onClick={props.handleClose}>
            Close
          </Button>
        </Tabs>
      </div>
      <div className="mt-2 grow mb-4">
        {
          {
            result: (
              // <div className="bg-gray-600 rounded-md m-2 p-2">
              //   <div className=" m-2 w-full xs:w-[300px] sm:w-[640px] hd:w-full wide:w-[1160px] ml-auto mr-auto ">
              <Summary data={props.data} />
              //   </div>
              // </div>
            ),
            config: <Config data={props.data} />,
            debug: (
              <Debugger data={props.parsed} team={props.data.char_names} searchable={msgs} />
            ),
            details: <Details data={props.data} />,
            share: <Share data={props.data} />,
          }[tabID]
        }
      </div>
      {tabID === "debug" ? (
        <div className="w-full pl-2 pr-2">
          <ButtonGroup fill>
            <Button
              onClick={() => setOptOpen(true)}
              icon="cog"
              intent="primary"
            >
              <Trans>viewer.debug_settings</Trans>
            </Button>
          </ButtonGroup>
        </div>
      ) : null}

      <Options {...optProps} />
      <SendToSim
        config={props.data.config_file}
        isOpen={open}
        onClose={() => setOpen(false)}
      />
    </div>
  );
}

type SendToSimProps = {
  config: string;
  isOpen: boolean;
  onClose: () => void;
};

function SendToSim({ config, isOpen, onClose }: SendToSimProps) {
  const dispatch = useAppDispatch();
  const [_, setLocation] = useLocation();
  const [keepExistingTeam, setKeepExistingTeam] = React.useState<boolean>(
    () => {
      const saved = localStorage.getItem(LOCALSTORAGE_KEY);
      if (saved === "true") {
        return true;
      }
      return false;
    }
  );

  const openInSim = () => {
    onClose();
    dispatch(updateCfg(config, keepExistingTeam));
    setLocation("/simulator");
  };

  const handleToggleSelected = () => {
    localStorage.setItem(LOCALSTORAGE_KEY, keepExistingTeam ? "false" : "true");
    setKeepExistingTeam(!keepExistingTeam);
  };

  return (
    <Dialog isOpen={isOpen} onClose={onClose}>
      <div className={Classes.DIALOG_BODY}>
        <Trans>viewer.load_this_configuration</Trans>
        <Callout intent="warning" className="mt-2">
          <Trans>viewer.this_will_overwrite</Trans>
        </Callout>
        <Checkbox
          label="Copy action list only (ignore character stats)"
          className="mt-2"
          checked={keepExistingTeam}
          onClick={handleToggleSelected}
        />
      </div>

      <div className={Classes.DIALOG_FOOTER}>
        <div className={Classes.DIALOG_FOOTER_ACTIONS}>
          <Button onClick={openInSim} intent="primary">
            <Trans>viewer.continue</Trans>
          </Button>
          <Button onClick={onClose}>
            <Trans>viewer.cancel</Trans>
          </Button>
        </div>
      </div>
    </Dialog>
  );
}

type ViewerProps = {
  data: SimResults;
  className?: string;
  handleClose: () => void;
};

const SAVED_DEBUG_KEY = "gcsim-debug-settings";

export function Viewer(props: ViewerProps) {
  const [selected, setSelected] = React.useState<string[]>(() => {
    const saved = localStorage.getItem(SAVED_DEBUG_KEY);
    if (saved) {
      const initialValue = JSON.parse(saved);
      return initialValue || defOpts;
    }
    return defOpts;
  });

  //string
  console.log(props.data);

  let parsed: DebugRow[];
  if (props.data.v2) {
    console.log("parsing as v2: " + props.data.debug);
    parsed = parseLogV2(
      props.data.active_char,
      props.data.char_names,
      props.data.debug,
      selected
    );
  } else {
    console.log("parsing as v1: " + props.data.debug);
    parsed = parseLog(
      props.data.active_char,
      props.data.char_names,
      props.data.debug,
      selected
    );
  }

  console.log(parsed);

  const handleSetSelected = (next: string[]) => {
    setSelected(next);
    localStorage.setItem(SAVED_DEBUG_KEY, JSON.stringify(next));
  };

  let viewProps = {
    classes: props.className,
    selected: selected,
    handleSetSelected: handleSetSelected,
    data: props.data,
    parsed: parsed,
    handleClose: props.handleClose,
  };

  return <ViewOnly {...viewProps} />;
}
