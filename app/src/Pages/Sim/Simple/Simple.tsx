import {
  Button,
  Callout,
  Card,
  Collapse,
  Intent,
  useHotkeys,
} from "@blueprintjs/core";
import React from "react";

import { SectionDivider } from "~src/Components/SectionDivider";
import { Viewport } from "~src/Components/Viewport";
import { RootState, useAppDispatch, useAppSelector } from "~src/store";
import { simActions } from "..";
import { ActionList, SimOptions } from "../Components";
import { SimProgress } from "../Components/SimProgress";
import { runSim } from "../exec";
import { Team } from "./Team";

export function Simple() {
  const { ready, workers, cfg, runState, showTips } = useAppSelector(
    (state: RootState) => {
      return {
        ready: state.sim.ready,
        workers: state.sim.workers,
        cfg: state.sim.cfg,
        runState: state.sim.run,
        showTips: state.sim.showTips,
      };
    }
  );
  const dispatch = useAppDispatch();

  const [open, setOpen] = React.useState<boolean>(false);
  const [showActionList, setShowActionList] = React.useState<boolean>(true);
  const [showOptions, setShowOptions] = React.useState<boolean>(false);

  const hotkeys = React.useMemo(
    () => [
      {
        combo: "Esc",
        global: true,
        label: "Exit edit",
        onKeyDown: () => {
          dispatch(simActions.editCharacter({ index: -1 }));
        },
      },
    ],
    []
  );
  useHotkeys(hotkeys);

  const run = () => {
    dispatch(runSim(cfg));
    setOpen(true);
  };

  const toggleTips = () => {
    dispatch(simActions.setShowTips(!showTips));
  };

  return (
    <Viewport className="flex flex-col gap-2">
      <div className="flex flex-col">
        {!showTips ? (
          <div className="ml-auto mr-2">
            <Button icon="help" onClick={toggleTips}>
              Give me back my tooltips!
            </Button>
          </div>
        ) : null}

        <Team />
        <SectionDivider>Action List</SectionDivider>
        <div className="ml-auto mr-2">
          <Button
            icon="edit"
            onClick={() => setShowActionList(!showActionList)}
          >
            {showActionList ? "Hide" : "Show"}
          </Button>
        </div>
        {showTips ? (
          <div className="pl-2 pr-2 pt-2">
            <Callout intent={Intent.PRIMARY} className="flex flex-col">
              <p>
                Enter the action list you wish to use in the textbox below. We
                have a vast collection of user submitted action lists on our{" "}
                <a href="https://discord.gg/W36ZwwhEaG" target="_blank">
                  Discord
                </a>{" "}
                if you need some ideas.
              </p>
              <p>
                Alternatively, check out our{" "}
                <a
                  href="https://docs.gcsim.app/guide/sequential_mode"
                  target="_blank"
                >
                  documentation
                </a>{" "}
                which includes a simple tutorial on how to write sequential
                lists (previously known as calc mode)
              </p>
              <div className="ml-auto">
                <Button small onClick={toggleTips}>
                  Hide all tips
                </Button>
              </div>
            </Callout>
          </div>
        ) : null}
        <Collapse
          isOpen={showActionList}
          keepChildrenMounted
          className="basis-full flex flex-col"
        >
          <ActionList
            cfg={cfg}
            onChange={(v) => dispatch(simActions.setCfg(v))}
          />
        </Collapse>
        <SectionDivider>Sim Options</SectionDivider>
        <div className="ml-auto mr-2">
          <Button icon="edit" onClick={() => setShowOptions(!showOptions)}>
            {showOptions ? "Hide" : "Show"}
          </Button>
        </div>
        <Collapse
          isOpen={showOptions}
          keepChildrenMounted
          className="basis-full flex flex-col"
        >
          <SimOptions />
        </Collapse>
      </div>
      <div className="sticky bottom-0 bg-bp-bg p-2 wide:ml-2 wide:mr-2 flex flex-row flex-wrap place-items-center gap-x-1 gap-y-1">
        <div className="basis-full wide:basis-0 flex-grow p-1">
          {`Workers available: ${ready}`}
        </div>
        <div className="basis-full wide:basis-1/3 p-1">
          <Button
            icon="play"
            fill
            intent="primary"
            onClick={run}
            disabled={ready < workers || runState.progress !== -1}
          >
            {ready < workers ? "Loading workers" : "Run"}
          </Button>
        </div>
      </div>
      <SimProgress isOpen={open} onClose={() => setOpen(false)} />
    </Viewport>
  );
}
