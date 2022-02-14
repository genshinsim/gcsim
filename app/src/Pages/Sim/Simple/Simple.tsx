import { Button, Callout, Card, Collapse, Intent } from "@blueprintjs/core";
import React from "react";
import { NumberInput } from "~src/Components/NumberInput";
import { SectionDivider } from "~src/Components/SectionDivider";
import { Viewport } from "~src/Components/Viewport";
import { RootState, useAppDispatch, useAppSelector } from "~src/store";
import { simActions } from "..";
import { ActionList } from "../Components";
import { SimProgress } from "../Components/SimProgress";
import { runSim } from "../exec";
import { Team } from "./Team";

export function Simple() {
  const { ready, workers } = useAppSelector((state: RootState) => {
    return {
      ready: state.sim.ready,
      workers: state.sim.workers,
    };
  });
  const dispatch = useAppDispatch();
  const [open, setOpen] = React.useState<boolean>(false);
  const [showActionList, setShowActionList] = React.useState<boolean>(true);
  const [showOptions, setShowOptions] = React.useState<boolean>(true);

  const run = () => {
    dispatch(runSim());
    setOpen(true);
  };
  return (
    <Viewport className="flex flex-col gap-2">
      <div className="flex flex-col">
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
        <div className="pl-2 pr-2 pt-2">
          <Callout intent={Intent.PRIMARY} className="flex flex-col">
            Enter action list here. For more detailed on action list, see ??
            <br />
            <div className="ml-auto">
              <Button small>Hide all tips</Button>
            </div>
          </Callout>
        </div>
        <Collapse
          isOpen={showActionList}
          keepChildrenMounted
          className="basis-full flex flex-col"
        >
          <ActionList />
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
          <Card className="m-2">
            <div className="w-full wide:basis-0 flex-grow p-2 text-center">
              <NumberInput
                label={`Workers (max available: ${ready})`}
                onChange={(v) => dispatch(simActions.setWorkers(v))}
                value={workers}
                min={1}
                max={ready}
                integerOnly
              />
            </div>
          </Card>
        </Collapse>
      </div>
      <div className="sticky bottom-0 bg-bp-bg p-2 wide:ml-2 wide:mr-2 flex flex-row flex-wrap place-items-center gap-x-1 gap-y-1">
        <div className="ml-auto basis-full wide:basis-1/3 p-1">
          <Button
            icon="play"
            fill
            intent="primary"
            onClick={run}
            disabled={ready < workers}
          >
            {ready < workers ? "Loading workers" : "Run"}
          </Button>
        </div>
      </div>
      <SimProgress isOpen={open} onClose={() => setOpen(false)} />
    </Viewport>
  );
}
