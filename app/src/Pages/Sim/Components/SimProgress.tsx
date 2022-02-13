import { Button, Dialog, ProgressBar } from "@blueprintjs/core";
import { useAppSelector, RootState, useAppDispatch } from "~src/store";

type Props = {
  isOpen: boolean;
  onClose: () => void;
};

export function SimProgress(props: Props) {
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
      <div className="p-2">
        {!done ? (
          <div className="flex flex-col gap-1">
            <div>Running sim with {workers} workers</div>
            <ProgressBar animate intent="primary" value={run.progress / 20} />
          </div>
        ) : (
          <div className="flex flex-col gap-1">
            <div>
              {" "}
              Simulation completed in {run.time}ms with average dps:{" "}
              {run.result}{" "}
            </div>
            <Button fill onClick={props.onClose}>
              done
            </Button>
          </div>
        )}
      </div>
    </Dialog>
  );
}
