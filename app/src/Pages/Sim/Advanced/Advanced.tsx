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

export function Advanced() {
  const { ready, workers, cfg, runState } = useAppSelector(
    (state: RootState) => {
      return {
        ready: state.sim.ready,
        workers: state.sim.workers,
        cfg: state.sim.advanced_cfg,
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
        <SectionDivider>Action List</SectionDivider>
        <ActionList
          cfg={cfg}
          onChange={(v) => dispatch(simActions.setAdvCfg(v))}
        />
        <SectionDivider>Helpers</SectionDivider>
        <div className="p-2">
          <Button disabled>Substat Helper</Button>
        </div>
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
