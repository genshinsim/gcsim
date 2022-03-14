import { Button, Classes, Dialog, ProgressBar } from "@blueprintjs/core";
import { RootState, useAppDispatch, useAppSelector } from "~src/store";
import { useLocation } from "wouter";
import { Trans, useTranslation } from "react-i18next";

type Props = {
  isOpen: boolean;
  onClose: () => void;
};

export function SimProgress(props: Props) {
  useTranslation()

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
  // console.log(run);

  return (
    <Dialog
      isOpen={props.isOpen}
      canEscapeKeyClose={done}
      canOutsideClickClose={done}
      onClose={props.onClose}
    >
      <div className="flex flex-col rounded-md">
        <div className="text-left text-lg bg-gray-600  rounded-t-md mb-4 pl-2 pb-2 pt-2">
          <Trans>components.running_simulation</Trans>
        </div>
        <div className="p-4 flex-grow">
          {!done ? (
            <div className="flex flex-col gap-1">
              <div><Trans>components.workers_pre</Trans>{workers}<Trans>components.workers_post</Trans></div>
              <ProgressBar animate intent="primary" value={run.progress / 20} />
            </div>
          ) : (
            <div className="flex flex-col gap-1">
              {run.err === "" ? (
                <div>
                  <Trans>components.result_pre</Trans>{run.time.toFixed(0)}<Trans>components.result_post</Trans>{run.result.toFixed(0)}
                </div>
              ) : (
                <div>
                  <Trans>components.simulation_exited_with</Trans>
                  <pre className="p-2 mt-2 whitespace-pre-wrap bg-gray-600 rounded-md">
                    {run.err}
                  </pre>
                </div>
              )}
            </div>
          )}
        </div>
        <div className={Classes.DIALOG_FOOTER}>
          <div className={Classes.DIALOG_FOOTER_ACTIONS}>
            <Button
              onClick={() => setLocation("/viewer")}
              disabled={!done || run.err !== ""}
              intent="success"
            >
              <Trans>components.see_results_in</Trans>
            </Button>
            <Button onClick={props.onClose} disabled={!done} intent="danger">
              <Trans>components.close</Trans>
            </Button>
          </div>
        </div>
      </div>
    </Dialog>
  );
}
