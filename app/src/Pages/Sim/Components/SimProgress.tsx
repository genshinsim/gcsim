import { Button, Classes, Dialog, ProgressBar } from "@blueprintjs/core";
import { Tooltip2 } from "@blueprintjs/popover2";
import { useAppSelector, RootState, useAppDispatch } from "~src/store";
import { useLocation } from "wouter";

type Props = {
  isOpen: boolean;
  onClose: () => void;
};

export function SimProgress(props: Props) {
  const [_, setLocation] = useLocation();
  const { run, workers } = useAppSelector((state: RootState) => {
    return {
      run: state.sim.run,
      workers: state.sim.workers,
    };
  });
  const dispatch = useAppDispatch();

  let done = run.progress === -1;
  //   console.log(done);

  return (
    <Dialog
      isOpen={props.isOpen}
      canEscapeKeyClose={done}
      canOutsideClickClose={done}
      onClose={props.onClose}
    >
      <div className="flex flex-col rounded-md">
        <div className="text-left text-lg bg-gray-600  rounded-t-md mb-4 pl-2 pb-2 pt-2">
          Running Simulation:
        </div>
        <div className="p-4 flex-grow">
          {!done ? (
            <div className="flex flex-col gap-1">
              <div>Running sim with {workers} workers</div>
              <ProgressBar animate intent="primary" value={run.progress / 20} />
            </div>
          ) : (
            <div className="flex flex-col gap-1">
              <div>
                Simulation completed in {run.time.toFixed(0)}ms with average
                dps: {run.result.toFixed(0)}{" "}
              </div>
            </div>
          )}
        </div>
        <div className={Classes.DIALOG_FOOTER}>
          <div className={Classes.DIALOG_FOOTER_ACTIONS}>
            <Button onClick={() => setLocation('/viewer')} disabled={!done} intent="success">
              See Results in Viewer
            </Button>
            <Button onClick={props.onClose} disabled={!done} intent="danger">
              Close
            </Button>
          </div>
        </div>
      </div>
    </Dialog>
  );
}
