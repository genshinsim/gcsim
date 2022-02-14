import { Button, Tab, Tabs } from "@blueprintjs/core";
import React from "react";
import { Config } from "./Config";
import { SimResults } from "./DataType";
import { Debugger } from "./DebugView";
import { Details } from "./Details";
import { Options, OptionsProp } from "./Options";
import { DebugRow, parseLog } from "./parse";
import Summary from "./Summary";
import Ajv from "ajv";
import schema from "./DataType.schema.json";
import Share, { ShareProps } from "./Share";
import { RootState, useAppSelector } from "~src/store";

const ajv = new Ajv();

type ViewerProps = {
  data: string;
  names?: string;
  handleClose: () => void;
};

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

function ViewOnly(props: ViewProps) {
  const [tabID, setTabID] = React.useState<string>("result");
  const [optOpen, setOptOpen] = React.useState<boolean>(false);
  const [shareOpen, setShareOpen] = React.useState<boolean>(false);

  const handleTabChange = (next: string) => {
    if (next === "settings") {
      setOptOpen(true);
      return;
    }
    if (next == "share") {
      setShareOpen(true);
      return;
    }
    setTabID(next);
  };

  const shareProps: ShareProps = {
    isOpen: shareOpen,
    handleClose: () => {
      setShareOpen(false);
    },
    data: props.data,
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

  return (
    <div
      className={props.classes + " p-4 rounded-lg bg-gray-800 flex flex-col"}
    >
      <div className="flex flex-row  bg-gray-800 w-full">
        <Tabs
          selectedTabId={tabID}
          onChange={handleTabChange}
          className="w-full"
        >
          <Tab id="result" title="Summary" className="focus:outline-none" />
          <Tab id="details" title="Details" className="focus:outline-none" />
          <Tab id="config" title="Config" className="focus:outline-none" />
          <Tab id="debug" title="Debug" className="focus:outline-none" />
          <Tab id="settings" title="Settings" className="focus:outline-none" />
          <Tab id="share" title="Share" className="focus:outline-none" />
          <Tabs.Expander />
          <Button icon="cross" intent="danger" onClick={props.handleClose} />
        </Tabs>
      </div>
      <div className="mt-2 grow">
        {
          {
            result: <Summary data={props.data} />,
            config: <Config data={props.data} />,
            debug: (
              <Debugger data={props.parsed} team={props.data.char_names} />
            ),
            details: <Details data={props.data} />,
          }[tabID]
        }
      </div>

      <Options {...optProps} />
      <Share {...shareProps} />
    </div>
  );
}

export function Viewer(props: ViewerProps) {
  const [selected, setSelected] = React.useState<string[]>(defOpts);
  const { simResults } = useAppSelector((state: RootState) => {
    return {
      simResults: state.sim.simResults,
    };
  });

  //string
  console.log(simResults);

  let data: SimResults = JSON.parse(props.data != '{}' ? props.data : simResults);
  const validate = ajv.compile(schema.definitions["*"]);
  const valid = validate(data);
  console.log("checking if data is valid: " + valid);

  if (!valid) {
    console.log(validate.errors);
    return (
      <div
        className={
          props.names +
          " p-4 rounded-lg bg-gray-800 flex flex-col w-full place-content-center items-center"
        }
      >
        <div className="mb-4 text-center">
          The data you have provided is not a valid format.{" "}
          <span className="font-bold">
            Please make sure you are using gcsim version 0.4.25 or higher.
          </span>
          <br />
          <br />
          Please click the close button and upload a valid file.
        </div>
        <div>
          <Button intent="danger" icon="cross" onClick={props.handleClose}>
            Click Here To Close
          </Button>
        </div>
        <div className="mt-8 rounded-md p-4 bg-gray-600">
          <p>
            If you think this error is invalid, please show the following
            message to the developers
          </p>
          <pre>{JSON.stringify(validate.errors, null, 2)}</pre>
        </div>
      </div>
    );
  }

  const parsed = parseLog(
    data.active_char,
    data.char_names,
    data.debug,
    selected
  );

  const handleSetSelected = (next: string[]) => {
    setSelected(next);
  };

  let viewProps = {
    classes: props.names,
    selected: selected,
    handleSetSelected: handleSetSelected,
    data: data,
    parsed: parsed,
    handleClose: props.handleClose,
  };

  return <ViewOnly {...viewProps} />;
}
